//go:build freebsd || openbsd || netbsd
// +build freebsd openbsd netbsd

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import "syscall"

func direntInode(dirent *syscall.Dirent) uint64 {
	return uint64(dirent.Fileno)
}
