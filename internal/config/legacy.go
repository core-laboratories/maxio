// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package config

// One time migration code section

// SetRegion - One time migration code needed, for migrating from older config to new for server Region.
func SetRegion(c Config, name string) {
	if name == "" {
		return
	}
	c[RegionSubSys][Default] = KVS{
		KV{
			Key:   RegionName,
			Value: name,
		},
	}
}
