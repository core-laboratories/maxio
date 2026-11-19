// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ioutil

import (
	"io"
)

// Discard is just like io.Discard without the io.ReaderFrom compatible
// implementation which is buggy on NUMA systems, we have to use a simpler
// io.Writer implementation alone avoids also unnecessary buffer copies,
// and as such incurred latencies.
var Discard io.Writer = discard{}

// discard is /dev/null for Golang.
type discard struct{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}

// DiscardReader discarded reader
func DiscardReader(r io.Reader) {
	Copy(Discard, r)
}
