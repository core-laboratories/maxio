// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package callhome

import "github.com/core-laboratories/maxio/internal/config"

var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	// HelpCallhome - provides help for callhome config
	HelpCallhome = config.HelpKVS{
		config.HelpKV{
			Key:         Enable,
			Type:        "on|off",
			Description: "set to enable callhome" + defaultHelpPostfix(Enable),
			Optional:    true,
		},
		config.HelpKV{
			Key:         Frequency,
			Type:        "duration",
			Description: "time duration between callhome cycles e.g. 24h" + defaultHelpPostfix(Frequency),
			Optional:    true,
		},
	}
)
