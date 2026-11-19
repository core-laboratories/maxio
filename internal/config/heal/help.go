// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package heal

import "github.com/core-laboratories/maxio/internal/config"

// Help template for caching feature.
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	// Help provides help for config values
	Help = config.HelpKVS{
		config.HelpKV{
			Key:         Bitrot,
			Description: `perform bitrot scan on drives when checking objects during scanner` + defaultHelpPostfix(Bitrot),
			Optional:    true,
			Type:        "on|off",
		},
		config.HelpKV{
			Key:         Sleep,
			Description: `maximum sleep duration between objects to slow down heal operation` + defaultHelpPostfix(Sleep),
			Optional:    true,
			Type:        "duration",
		},
		config.HelpKV{
			Key:         IOCount,
			Description: `maximum IO requests allowed between objects to slow down heal operation` + defaultHelpPostfix(IOCount),
			Optional:    true,
			Type:        "int",
		},
		config.HelpKV{
			Key:         DriveWorkers,
			Description: `the number of workers per drive to heal a new disk replacement` + defaultHelpPostfix(DriveWorkers),
			Optional:    true,
			Type:        "int",
		},
	}
)
