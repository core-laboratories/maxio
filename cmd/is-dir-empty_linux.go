//go:build linux && !appengine
// +build linux,!appengine

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"syscall"
)

// Returns true if no error and there is no object or prefix inside this directory
func isDirEmpty(dirname string, legacy bool) bool {
	if legacy {
		// On filesystems such as btrfs, nfs this is not true, so fallback
		// to performing readdir() instead.
		entries, err := readDirN(dirname, 1)
		if err != nil {
			return false
		}
		return len(entries) == 0
	}
	var stat syscall.Stat_t
	if err := syscall.Stat(dirname, &stat); err != nil {
		return false
	}
	return stat.Mode&syscall.S_IFMT == syscall.S_IFDIR && stat.Nlink == 2
}
