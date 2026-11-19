// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package hash

import (
	"bytes"
	"errors"
	"hash"
	"io"

	"github.com/core-laboratories/maxio/internal/ioutil"
)

// Checker allows to verify the checksum of a reader.
type Checker struct {
	c io.Closer
	r io.Reader
	h hash.Hash

	want []byte
}

// NewChecker ensures that content with the specified length is read from rc.
// Calling Close on this will close upstream.
func NewChecker(rc io.ReadCloser, h hash.Hash, wantSum []byte, length int64) *Checker {
	return &Checker{c: rc, r: ioutil.HardLimitReader(rc, length), h: h, want: wantSum}
}

// Read satisfies io.Reader
func (c Checker) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	if n > 0 {
		c.h.Write(p[:n])
	}
	if errors.Is(err, io.EOF) {
		got := c.h.Sum(nil)
		if !bytes.Equal(got, c.want) {
			return n, ErrInvalidChecksum
		}
		return n, err
	}
	return n, err
}

// Close satisfies io.Closer
func (c Checker) Close() error {
	err := c.c.Close()
	if err == nil {
		got := c.h.Sum(nil)
		if !bytes.Equal(got, c.want) {
			return ErrInvalidChecksum
		}
	}
	return err
}
