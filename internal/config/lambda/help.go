// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lambda

import (
	"github.com/core-laboratories/maxio/internal/config"
	"github.com/core-laboratories/maxio/internal/event/target"
)

// Help template inputs for all lambda targets
var (
	HelpWebhook = config.HelpKVS{
		config.HelpKV{
			Key:         target.WebhookEndpoint,
			Description: "webhook server endpoint e.g. http://localhost:8080/minio/lambda",
			Type:        "url",
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         target.WebhookAuthToken,
			Description: "opaque string or JWT authorization token",
			Optional:    true,
			Type:        "string",
			Sensitive:   true,
			Secret:      true,
		},
		config.HelpKV{
			Key:         config.Comment,
			Description: config.DefaultComment,
			Optional:    true,
			Type:        "sentence",
		},
		config.HelpKV{
			Key:         target.WebhookClientCert,
			Description: "client cert for Webhook mTLS auth",
			Optional:    true,
			Type:        "string",
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         target.WebhookClientKey,
			Description: "client cert key for Webhook mTLS auth",
			Optional:    true,
			Type:        "string",
			Sensitive:   true,
		},
	}
)
