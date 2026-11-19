//go:build !linux
// +build !linux

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

// isDirEmpty - returns true if there is no error and no object and prefix inside this directory
func isDirEmpty(dirname string, _ bool) bool {
	entries, err := readDirN(dirname, 1)
	if err != nil {
		return false
	}
	return len(entries) == 0
}
