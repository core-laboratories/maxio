// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ioutil

import (
	"io"
	"os"

	"github.com/core-laboratories/maxio/internal/lock"
)

// AppendFile - appends the file "src" to the file "dst"
func AppendFile(dst string, src string, osync bool) error {
	appendFile, err := lock.Open(dst, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer appendFile.Close()

	srcFile, err := lock.Open(src, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	_, err = io.Copy(appendFile, srcFile)
	return err
}
