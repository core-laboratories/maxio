// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package batch

import "github.com/core-laboratories/maxio/internal/config"

// Help template for batch feature.
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	// Help provides help for config values
	Help = config.HelpKVS{
		config.HelpKV{
			Key:         ReplicationWorkersWait,
			Description: `maximum sleep duration between objects to slow down batch replication operation` + defaultHelpPostfix(ReplicationWorkersWait),
			Optional:    true,
			Type:        "duration",
		},
		config.HelpKV{
			Key:         KeyRotationWorkersWait,
			Description: `maximum sleep duration between objects to slow down batch keyrotation operation` + defaultHelpPostfix(KeyRotationWorkersWait),
			Optional:    true,
			Type:        "duration",
		},
		config.HelpKV{
			Key:         ExpirationWorkersWait,
			Description: "maximum sleep duration between objects to slow down batch expiration operation" + defaultHelpPostfix(ExpirationWorkersWait),
			Optional:    true,
			Type:        "duration",
		},
	}
)
