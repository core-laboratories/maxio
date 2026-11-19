//go:build !linux && !windows
// +build !linux,!windows

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package mountinfo

// CheckCrossDevice - check if any input path has multiple sub-mounts.
// this is a dummy function and returns nil for now.
func CheckCrossDevice(paths []string) error {
	return nil
}

// IsLikelyMountPoint determines if a directory is a mountpoint.
func IsLikelyMountPoint(file string) bool {
	return false
}
