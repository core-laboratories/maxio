// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package storageclass

import "github.com/core-laboratories/maxio/internal/config"

// Help template for storageclass feature.
var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	Help = config.HelpKVS{
		config.HelpKV{
			Key:         ClassStandard,
			Description: `set the parity count for default standard storage class` + defaultHelpPostfix(ClassStandard),
			Optional:    true,
			Type:        "string",
		},
		config.HelpKV{
			Key:         ClassRRS,
			Description: `set the parity count for reduced redundancy storage class` + defaultHelpPostfix(ClassRRS),
			Optional:    true,
			Type:        "string",
		},
		config.HelpKV{
			Key:         Optimize,
			Description: `optimize parity calculation for standard storage class, set 'capacity' for capacity optimized (no additional parity)` + defaultHelpPostfix(Optimize),
			Optional:    true,
			Type:        "string",
		},
		config.HelpKV{
			Key:         config.Comment,
			Description: config.DefaultComment,
			Optional:    true,
			Type:        "sentence",
		},
	}
)
