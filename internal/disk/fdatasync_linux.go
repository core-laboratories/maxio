//go:build linux
// +build linux

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package disk

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// Fdatasync - fdatasync() is similar to fsync(), but does not flush modified metadata
// unless that metadata is needed in order to allow a subsequent data retrieval
// to  be  correctly  handled.   For example, changes to st_atime or st_mtime
// (respectively, time of last access and time of last modification; see inode(7))
// do not require flushing because they are not necessary for a subsequent data
// read to be handled correctly. On the other hand, a change to the file size
// (st_size, as made by say ftruncate(2)), would require a metadata flush.
//
// The aim of fdatasync() is to reduce disk activity for applications that
// do not require all metadata to be synchronized with the disk.
func Fdatasync(f *os.File) error {
	return syscall.Fdatasync(int(f.Fd()))
}

// FadviseDontNeed invalidates page-cache
func FadviseDontNeed(f *os.File) error {
	return unix.Fadvise(int(f.Fd()), 0, 0, unix.FADV_DONTNEED)
}
