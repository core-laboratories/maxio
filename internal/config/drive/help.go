// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package drive

import "github.com/core-laboratories/maxio/internal/config"

var (
	// MaxTimeout is the max timeout for drive
	MaxTimeout = "max_timeout"

	// HelpDrive is help for drive
	HelpDrive = config.HelpKVS{
		config.HelpKV{
			Key:         MaxTimeout,
			Type:        "string",
			Description: "set per call max_timeout for the drive, defaults to 30 seconds",
			Optional:    true,
		},
	}
)
