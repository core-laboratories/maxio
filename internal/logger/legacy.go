// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package logger

import (
	"github.com/core-laboratories/maxio/internal/config"
	"github.com/core-laboratories/maxio/internal/logger/target/http"
)

// Legacy envs
const (
	legacyEnvAuditLoggerHTTPEndpoint = "MINIO_AUDIT_LOGGER_HTTP_ENDPOINT"
	legacyEnvLoggerHTTPEndpoint      = "MINIO_LOGGER_HTTP_ENDPOINT"
)

// SetLoggerHTTPAudit - helper for migrating older config to newer KV format.
func SetLoggerHTTPAudit(scfg config.Config, k string, args http.Config) {
	if !args.Enabled {
		// Do not enable audit targets, if not enabled
		return
	}
	scfg[config.AuditWebhookSubSys][k] = config.KVS{
		config.KV{
			Key:   config.Enable,
			Value: config.EnableOn,
		},
		config.KV{
			Key:   Endpoint,
			Value: args.Endpoint.String(),
		},
		config.KV{
			Key:   AuthToken,
			Value: args.AuthToken,
		},
	}
}

// SetLoggerHTTP helper for migrating older config to newer KV format.
func SetLoggerHTTP(scfg config.Config, k string, args http.Config) {
	if !args.Enabled {
		// Do not enable logger http targets, if not enabled
		return
	}

	scfg[config.LoggerWebhookSubSys][k] = config.KVS{
		config.KV{
			Key:   config.Enable,
			Value: config.EnableOn,
		},
		config.KV{
			Key:   Endpoint,
			Value: args.Endpoint.String(),
		},
		config.KV{
			Key:   AuthToken,
			Value: args.AuthToken,
		},
	}
}
