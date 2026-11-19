// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lifecycle

import "testing"

func Test_NoncurrentVersionsExpiration_Validation(t *testing.T) {
	testcases := []struct {
		n   NoncurrentVersionExpiration
		err error
	}{
		{
			n: NoncurrentVersionExpiration{
				NoncurrentDays:          0,
				NewerNoncurrentVersions: 0,
				set:                     true,
			},
			err: errXMLNotWellFormed,
		},
		{
			n: NoncurrentVersionExpiration{
				NoncurrentDays:          90,
				NewerNoncurrentVersions: 0,
				set:                     true,
			},
			err: nil,
		},
		{
			n: NoncurrentVersionExpiration{
				NoncurrentDays:          90,
				NewerNoncurrentVersions: 2,
				set:                     true,
			},
			err: nil,
		},
		{
			n: NoncurrentVersionExpiration{
				NoncurrentDays: -1,
				set:            true,
			},
			err: errXMLNotWellFormed,
		},
		{
			n: NoncurrentVersionExpiration{
				NoncurrentDays:          90,
				NewerNoncurrentVersions: -2,
				set:                     true,
			},
			err: errXMLNotWellFormed,
		},
		// MinIO extension: supports zero NoncurrentDays when NewerNoncurrentVersions > 0
		{
			n: NoncurrentVersionExpiration{
				NoncurrentDays:          0,
				NewerNoncurrentVersions: 5,
				set:                     true,
			},
			err: nil,
		},
	}

	for i, tc := range testcases {
		if got := tc.n.Validate(); got != tc.err {
			t.Fatalf("%d: expected %v but got %v", i+1, tc.err, got)
		}
	}
}
