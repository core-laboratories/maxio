// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

// Options for readDir function call
type readDirOpts struct {
	// The maximum number of entries to return
	count int
	// Follow directory symlink
	followDirSymlink bool
}

// Return all the entries at the directory dirPath.
func readDir(dirPath string) (entries []string, err error) {
	return readDirWithOpts(dirPath, readDirOpts{count: -1})
}

// Return up to count entries at the directory dirPath.
func readDirN(dirPath string, count int) (entries []string, err error) {
	return readDirWithOpts(dirPath, readDirOpts{count: count})
}
