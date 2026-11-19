// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package scanner

import "github.com/core-laboratories/maxio/internal/config"

var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	// Help provides help for config values
	Help = config.HelpKVS{
		config.HelpKV{
			Key:         Speed,
			Description: `customize scanner speed (default|slowest|slow|fast|fastest)` + defaultHelpPostfix(Speed),
			Optional:    true,
			Type:        "string",
		},
		config.HelpKV{
			Key:         ExcessVersions,
			Description: `alert per object beyond this many versions` + defaultHelpPostfix(ExcessVersions),
			Optional:    true,
			Type:        "int",
		},
		config.HelpKV{
			Key:         ExcessFolders,
			Description: `alert beyond this many sub-folders per folder in an erasure set` + defaultHelpPostfix(ExcessFolders),
			Optional:    true,
			Type:        "int",
		},
	}
)
