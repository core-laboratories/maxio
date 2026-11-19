// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package hash

import (
	"errors"
	"fmt"
)

// SHA256Mismatch - when content sha256 does not match with what was sent from client.
type SHA256Mismatch struct {
	ExpectedSHA256   string
	CalculatedSHA256 string
}

func (e SHA256Mismatch) Error() string {
	return "Bad sha256: Expected " + e.ExpectedSHA256 + " does not match calculated " + e.CalculatedSHA256
}

// BadDigest - Content-MD5 you specified did not match what we received.
type BadDigest struct {
	ExpectedMD5   string
	CalculatedMD5 string
}

func (e BadDigest) Error() string {
	return "Bad digest: Expected " + e.ExpectedMD5 + " does not match calculated " + e.CalculatedMD5
}

// SizeTooSmall reader size too small
type SizeTooSmall struct {
	Want int64
	Got  int64
}

func (e SizeTooSmall) Error() string {
	return fmt.Sprintf("Size small: got %d, want %d", e.Got, e.Want)
}

// SizeTooLarge reader size too large
type SizeTooLarge struct {
	Want int64
	Got  int64
}

func (e SizeTooLarge) Error() string {
	return fmt.Sprintf("Size large: got %d, want %d", e.Got, e.Want)
}

// SizeMismatch error size mismatch
type SizeMismatch struct {
	Want int64
	Got  int64
}

func (e SizeMismatch) Error() string {
	return fmt.Sprintf("Size mismatch: got %d, want %d", e.Got, e.Want)
}

// ChecksumMismatch - when content checksum does not match with what was sent from client.
type ChecksumMismatch struct {
	Want string
	Got  string
}

func (e ChecksumMismatch) Error() string {
	return "Bad checksum: Want " + e.Want + " does not match calculated " + e.Got
}

// IsChecksumMismatch matches if 'err' is hash.ChecksumMismatch
func IsChecksumMismatch(err error) bool {
	var herr ChecksumMismatch
	return errors.As(err, &herr)
}
