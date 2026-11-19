// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package etcd

import "github.com/core-laboratories/maxio/internal/config"

// etcd config documented in default config
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	Help = config.HelpKVS{
		config.HelpKV{
			Key:         Endpoints,
			Description: `comma separated list of etcd endpoints` + defaultHelpPostfix(Endpoints),
			Type:        "csv",
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         PathPrefix,
			Description: `namespace prefix to isolate tenants` + defaultHelpPostfix(PathPrefix),
			Optional:    true,
			Type:        "path",
		},
		config.HelpKV{
			Key:         CoreDNSPath,
			Description: `shared bucket DNS records` + defaultHelpPostfix(CoreDNSPath),
			Optional:    true,
			Type:        "path",
		},
		config.HelpKV{
			Key:         ClientCert,
			Description: `client cert for mTLS authentication` + defaultHelpPostfix(ClientCert),
			Optional:    true,
			Type:        "path",
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         ClientCertKey,
			Description: `client cert key for mTLS authentication` + defaultHelpPostfix(ClientCertKey),
			Optional:    true,
			Type:        "path",
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         config.Comment,
			Description: config.DefaultComment,
			Optional:    true,
			Type:        "sentence",
		},
	}
)
