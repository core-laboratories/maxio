// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"os"
	"runtime"
	"syscall"
	"testing"
)

func TestSysErrors(t *testing.T) {
	pathErr := &os.PathError{Err: syscall.ENAMETOOLONG}
	ok := isSysErrTooLong(pathErr)
	if !ok {
		t.Fatalf("Unexpected error expecting %s", syscall.ENAMETOOLONG)
	}
	pathErr = &os.PathError{Err: syscall.ENOTDIR}
	ok = isSysErrNotDir(pathErr)
	if !ok {
		t.Fatalf("Unexpected error expecting %s", syscall.ENOTDIR)
	}
	if runtime.GOOS != globalWindowsOSName {
		pathErr = &os.PathError{Err: syscall.ENOTEMPTY}
		ok = isSysErrNotEmpty(pathErr)
		if !ok {
			t.Fatalf("Unexpected error expecting %s", syscall.ENOTEMPTY)
		}
	} else {
		pathErr = &os.PathError{Err: syscall.Errno(0x91)}
		ok = isSysErrNotEmpty(pathErr)
		if !ok {
			t.Fatal("Unexpected error expecting 0x91")
		}
	}
	if runtime.GOOS == globalWindowsOSName {
		pathErr = &os.PathError{Err: syscall.Errno(0x03)}
		ok = isSysErrPathNotFound(pathErr)
		if !ok {
			t.Fatal("Unexpected error expecting 0x03")
		}
	}
}
