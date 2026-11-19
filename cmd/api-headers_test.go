// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"testing"
)

func TestNewRequestID(t *testing.T) {
	// Ensure that it returns an alphanumeric result of length 16.
	id := mustGetRequestID(UTCNow())

	if len(id) != 16 {
		t.Fail()
	}

	var e rune
	for _, char := range id {
		e = char

		// Ensure that it is alphanumeric, in this case, between 0-9 and A-Z.
		isAlnum := ('0' <= e && e <= '9') || ('A' <= e && e <= 'Z')
		if !isAlnum {
			t.Fail()
		}
	}
}
