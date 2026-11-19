// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/core-laboratories/maxio/internal/config"
	"github.com/core-laboratories/maxio/internal/config/compress"
	xldap "github.com/core-laboratories/maxio/internal/config/identity/ldap"
	"github.com/core-laboratories/maxio/internal/config/identity/openid"
	"github.com/core-laboratories/maxio/internal/config/notify"
	"github.com/core-laboratories/maxio/internal/config/policy/opa"
	"github.com/core-laboratories/maxio/internal/config/storageclass"
	"github.com/core-laboratories/maxio/internal/event/target"
	"github.com/core-laboratories/maxio/internal/logger"
	xnet "github.com/minio/pkg/v3/net"
	"github.com/minio/pkg/v3/quick"
)

// Save config file to corresponding backend
func Save(configFile string, data any) error {
	return quick.SaveConfig(data, configFile, globalEtcdClient)
}

// Load config from backend
func Load(configFile string, data any) (quick.Config, error) {
	return quick.LoadConfig(configFile, globalEtcdClient, data)
}

func readConfigWithoutMigrate(ctx context.Context, objAPI ObjectLayer) (config.Config, error) {
	// Construct path to config.json for the given bucket.
	configFile := path.Join(minioConfigPrefix, minioConfigFile)

	configFiles := []string{
		getConfigFile(),
		getConfigFile() + ".deprecated",
		configFile,
	}

	newServerCfg := func() (config.Config, error) {
		// Initialize server config.
		srvCfg := newServerConfig()

		return srvCfg, saveServerConfig(ctx, objAPI, srvCfg)
	}

	var data []byte
	var err error

	cfg := &serverConfigV33{}
	for _, cfgFile := range configFiles {
		if _, err = Load(cfgFile, cfg); err != nil {
			if !osIsNotExist(err) && !osIsPermission(err) {
				return nil, err
			}
			continue
		}
		data, _ = json.Marshal(cfg)
		break
	}
	if osIsPermission(err) {
		logger.Info("Older config found but is not readable %s, proceeding to read config from other places", err)
	}
	if osIsNotExist(err) || osIsPermission(err) || len(data) == 0 {
		data, err = readConfig(GlobalContext, objAPI, configFile)
		if err != nil {
			// when config.json is not found, then we freshly initialize.
			if errors.Is(err, errConfigNotFound) {
				return newServerCfg()
			}
			return nil, err
		}

		data, err = decryptData(data, configFile)
		if err != nil {
			return nil, err
		}

		newCfg, err := readServerConfig(GlobalContext, objAPI, data)
		if err == nil {
			return newCfg, nil
		}

		// Read older `.minio.sys/config/config.json`, if not
		// possible just fail.
		if err = json.Unmarshal(data, cfg); err != nil {
			// Unable to parse old JSON simply re-initialize a new one.
			return newServerCfg()
		}
	}

	if !globalCredViaEnv && cfg.Credential.IsValid() {
		// Preserve older credential if we do not have
		// root credentials set via environment variable.
		globalActiveCred = cfg.Credential
	}

	// Init compression config. For future migration, Compression config needs to be copied over from previous version.
	switch cfg.Version {
	case "29":
		// V29 -> V30
		cfg.Compression.Enabled = false
		cfg.Compression.Extensions = strings.Split(compress.DefaultExtensions, config.ValueSeparator)
		cfg.Compression.MimeTypes = strings.Split(compress.DefaultMimeTypes, config.ValueSeparator)
	case "30":
		// V30 -> V31
		cfg.OpenID = openid.Config{}
		cfg.Policy.OPA = opa.Args{
			URL:       &xnet.URL{},
			AuthToken: "",
		}
	case "31":
		// V31 -> V32
		cfg.Notify.NSQ = make(map[string]target.NSQArgs)
		cfg.Notify.NSQ["1"] = target.NSQArgs{}
	}

	// Move to latest.
	cfg.Version = "33"

	newCfg := newServerConfig()

	config.SetRegion(newCfg, cfg.Region)
	storageclass.SetStorageClass(newCfg, cfg.StorageClass)

	for k, loggerArgs := range cfg.Logger.HTTP {
		logger.SetLoggerHTTP(newCfg, k, loggerArgs)
	}
	for k, auditArgs := range cfg.Logger.AuditWebhook {
		logger.SetLoggerHTTPAudit(newCfg, k, auditArgs)
	}

	xldap.SetIdentityLDAP(newCfg, cfg.LDAPServerConfig)
	opa.SetPolicyOPAConfig(newCfg, cfg.Policy.OPA)
	compress.SetCompressionConfig(newCfg, cfg.Compression)

	for k, args := range cfg.Notify.AMQP {
		notify.SetNotifyAMQP(newCfg, k, args)
	}
	for k, args := range cfg.Notify.Elasticsearch {
		notify.SetNotifyES(newCfg, k, args)
	}
	for k, args := range cfg.Notify.Kafka {
		notify.SetNotifyKafka(newCfg, k, args)
	}
	for k, args := range cfg.Notify.MQTT {
		notify.SetNotifyMQTT(newCfg, k, args)
	}
	for k, args := range cfg.Notify.MySQL {
		notify.SetNotifyMySQL(newCfg, k, args)
	}
	for k, args := range cfg.Notify.NATS {
		notify.SetNotifyNATS(newCfg, k, args)
	}
	for k, args := range cfg.Notify.NSQ {
		notify.SetNotifyNSQ(newCfg, k, args)
	}
	for k, args := range cfg.Notify.PostgreSQL {
		notify.SetNotifyPostgres(newCfg, k, args)
	}
	for k, args := range cfg.Notify.Redis {
		notify.SetNotifyRedis(newCfg, k, args)
	}
	for k, args := range cfg.Notify.Webhook {
		notify.SetNotifyWebhook(newCfg, k, args)
	}

	return newCfg, nil
}
