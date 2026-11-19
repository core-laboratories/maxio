// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package plugin

import "github.com/core-laboratories/maxio/internal/config"

// Help template for Access Management Plugin policy feature.
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	Help = config.HelpKVS{
		config.HelpKV{
			Key:         URL,
			Description: `plugin hook endpoint (HTTP(S)) e.g. "http://localhost:8181/v1/data/httpapi/authz/allow"` + defaultHelpPostfix(URL),
			Type:        "url",
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         AuthToken,
			Description: "authorization header for plugin hook endpoint" + defaultHelpPostfix(AuthToken),
			Optional:    true,
			Type:        "string",
			Sensitive:   true,
			Secret:      true,
		},
		config.HelpKV{
			Key:         EnableHTTP2,
			Description: "Enable experimental HTTP2 support to connect to plugin service" + defaultHelpPostfix(EnableHTTP2),
			Optional:    true,
			Type:        "bool",
		},
		config.HelpKV{
			Key:         config.Comment,
			Description: config.DefaultComment,
			Optional:    true,
			Type:        "sentence",
		},
	}
)
