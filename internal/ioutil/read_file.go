// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ioutil

import (
	"io"
	"io/fs"
	"os"

	"github.com/core-laboratories/maxio/internal/disk"
)

var (
	// OpenFileDirectIO allows overriding default function.
	OpenFileDirectIO = disk.OpenFileDirectIO
	// OsOpen allows overriding default function.
	OsOpen = os.Open
	// OsOpenFile allows overriding default function.
	OsOpenFile = os.OpenFile
)

// ReadFileWithFileInfo reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func ReadFileWithFileInfo(name string) ([]byte, fs.FileInfo, error) {
	f, err := OsOpenFile(name, readMode, 0o666)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}

	dst := make([]byte, st.Size())
	_, err = io.ReadFull(f, dst)
	return dst, st, err
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
//
// passes NOATIME flag for reads on Unix systems to avoid atime updates.
func ReadFile(name string) ([]byte, error) {
	// Don't wrap with un-needed buffer.
	// Don't use os.ReadFile, since it doesn't pass NO_ATIME when present.
	f, err := OsOpenFile(name, readMode, 0o666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return io.ReadAll(f)
	}
	dst := make([]byte, st.Size())
	_, err = io.ReadFull(f, dst)
	return dst, err
}
