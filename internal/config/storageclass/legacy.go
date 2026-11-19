// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package storageclass

import (
	"github.com/core-laboratories/maxio/internal/config"
)

// SetStorageClass - One time migration code needed, for migrating from older config to new for StorageClass.
func SetStorageClass(s config.Config, cfg Config) {
	if len(cfg.Standard.String()) == 0 && len(cfg.RRS.String()) == 0 {
		// Do not enable storage-class if no settings found.
		return
	}
	s[config.StorageClassSubSys][config.Default] = config.KVS{
		config.KV{
			Key:   ClassStandard,
			Value: cfg.Standard.String(),
		},
		config.KV{
			Key:   ClassRRS,
			Value: cfg.RRS.String(),
		},
	}
}
