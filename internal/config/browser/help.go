// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package browser

import "github.com/core-laboratories/maxio/internal/config"

// Help template for browser feature.
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	Help = config.HelpKVS{
		config.HelpKV{
			Key:         browserCSPPolicy,
			Description: `set Content-Security-Policy response header value` + defaultHelpPostfix(browserCSPPolicy),
			Optional:    true,
			Type:        "string",
		},
		config.HelpKV{
			Key:         browserHSTSSeconds,
			Description: `set Strict-Transport-Security 'max-age' amount of seconds value` + defaultHelpPostfix(browserHSTSSeconds),
			Optional:    true,
			Type:        "number",
		},
		config.HelpKV{
			Key:         browserHSTSIncludeSubdomains,
			Description: `turn 'on' to set Strict-Transport-Security 'includeSubDomains' directive` + defaultHelpPostfix(browserHSTSIncludeSubdomains),
			Optional:    true,
			Type:        "boolean",
		},
		config.HelpKV{
			Key:         browserHSTSPreload,
			Description: `turn 'on' to set Strict-Transport-Security 'preload' directive` + defaultHelpPostfix(browserHSTSPreload),
			Optional:    true,
			Type:        "boolean",
		},
		config.HelpKV{
			Key:         browserReferrerPolicy,
			Description: `set Referrer-Policy response header value` + defaultHelpPostfix(browserReferrerPolicy),
			Optional:    true,
			Type:        "string",
		},
	}
)
