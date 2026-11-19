// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"github.com/minio/madmin-go/v3"
)

// GetTotalCapacity gets the total capacity in the cluster.
func GetTotalCapacity(diskInfo []madmin.Disk) (capacity uint64) {
	for _, disk := range diskInfo {
		capacity += disk.TotalSpace
	}
	return capacity
}

// GetTotalUsableCapacity gets the total usable capacity in the cluster.
func GetTotalUsableCapacity(diskInfo []madmin.Disk, s StorageInfo) (capacity uint64) {
	for _, disk := range diskInfo {
		// Ignore invalid.
		if disk.PoolIndex < 0 || len(s.Backend.StandardSCData) <= disk.PoolIndex {
			// https://github.com/core-laboratories/maxio/issues/16500
			continue
		}
		// Ignore parity disks
		if disk.DiskIndex < s.Backend.StandardSCData[disk.PoolIndex] {
			capacity += disk.TotalSpace
		}
	}
	return capacity
}

// GetTotalCapacityFree gets the total capacity free in the cluster.
func GetTotalCapacityFree(diskInfo []madmin.Disk) (capacity uint64) {
	for _, d := range diskInfo {
		capacity += d.AvailableSpace
	}
	return capacity
}

// GetTotalUsableCapacityFree gets the total usable capacity free in the cluster.
func GetTotalUsableCapacityFree(diskInfo []madmin.Disk, s StorageInfo) (capacity uint64) {
	for _, disk := range diskInfo {
		// Ignore invalid.
		if disk.PoolIndex < 0 || len(s.Backend.StandardSCData) <= disk.PoolIndex {
			// https://github.com/core-laboratories/maxio/issues/16500
			continue
		}
		// Ignore parity disks
		if disk.DiskIndex < s.Backend.StandardSCData[disk.PoolIndex] {
			capacity += disk.AvailableSpace
		}
	}
	return capacity
}
