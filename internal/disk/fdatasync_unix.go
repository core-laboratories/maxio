//go:build freebsd || netbsd || openbsd || darwin
// +build freebsd netbsd openbsd darwin

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
)

// Fdatasync is fsync on freebsd/darwin
func Fdatasync(f *os.File) error {
	return syscall.Fsync(int(f.Fd()))
}

// FadviseDontNeed is a no-op
func FadviseDontNeed(f *os.File) error {
	return nil
}
