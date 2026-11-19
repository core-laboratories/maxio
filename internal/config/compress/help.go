// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package compress

import "github.com/core-laboratories/maxio/internal/config"

// Help template for compress feature.
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	Help = config.HelpKVS{
		config.HelpKV{
			Key:         config.Enable,
			Description: "Enable or disable object compression",
			Type:        "on|off",
			Optional:    true,
			Sensitive:   false,
		},
		config.HelpKV{
			Key:         Extensions,
			Description: `comma separated file extensions` + defaultHelpPostfix(Extensions),
			Optional:    true,
			Type:        "csv",
		},
		config.HelpKV{
			Key:         MimeTypes,
			Description: `comma separated wildcard mime-types` + defaultHelpPostfix(MimeTypes),
			Optional:    true,
			Type:        "csv",
		},
		config.HelpKV{
			Key:         AllowEncrypted,
			Description: `enable 'encryption' along with compression`,
			Optional:    true,
			Type:        "on|off",
		},
		config.HelpKV{
			Key:         config.Comment,
			Description: config.DefaultComment,
			Optional:    true,
			Type:        "sentence",
		},
	}
)
