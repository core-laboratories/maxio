//go:build darwin

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package disk

import (
	"os"

	"github.com/ncw/directio"
	"golang.org/x/sys/unix"
)

// ODirectPlatform indicates if the platform supports O_DIRECT
const ODirectPlatform = true

// OpenFileDirectIO - bypass kernel cache.
func OpenFileDirectIO(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	return directio.OpenFile(filePath, flag, perm)
}

// DisableDirectIO - disables directio mode.
func DisableDirectIO(f *os.File) error {
	fd := f.Fd()
	_, err := unix.FcntlInt(fd, unix.F_NOCACHE, 0)
	return err
}

// AlignedBlock - pass through to directio implementation.
func AlignedBlock(blockSize int) []byte {
	return directio.AlignedBlock(blockSize)
}
