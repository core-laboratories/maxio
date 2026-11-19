// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lifecycle

import (
	"encoding/xml"
	"testing"
)

func TestTransitionUnmarshalXML(t *testing.T) {
	trTests := []struct {
		input string
		err   error
	}{
		{
			input: `<Transition>
			<Days>0</Days>
			<StorageClass>S3TIER-1</StorageClass>
		  </Transition>`,
			err: nil,
		},
		{
			input: `<Transition>
			<Days>1</Days>
			<Date>2021-01-01T00:00:00Z</Date>
			<StorageClass>S3TIER-1</StorageClass>
		  </Transition>`,
			err: errTransitionInvalid,
		},
		{
			input: `<Transition>
			<Days>1</Days>
		  </Transition>`,
			err: errXMLNotWellFormed,
		},
	}

	for i, tc := range trTests {
		var tr Transition
		err := xml.Unmarshal([]byte(tc.input), &tr)
		if err != nil {
			t.Fatalf("%d: xml unmarshal failed with %v", i+1, err)
		}
		if err = tr.Validate(); err != tc.err {
			t.Fatalf("%d: Invalid transition %v: err %v", i+1, tr, err)
		}
	}

	ntrTests := []struct {
		input string
		err   error
	}{
		{
			input: `<NoncurrentVersionTransition>
			<NoncurrentDays>0</NoncurrentDays>
			<StorageClass>S3TIER-1</StorageClass>
		  </NoncurrentVersionTransition>`,
			err: nil,
		},
		{
			input: `<NoncurrentVersionTransition>
			<Days>1</Days>
		  </NoncurrentVersionTransition>`,
			err: errXMLNotWellFormed,
		},
	}

	for i, tc := range ntrTests {
		var ntr NoncurrentVersionTransition
		err := xml.Unmarshal([]byte(tc.input), &ntr)
		if err != nil {
			t.Fatalf("%d: xml unmarshal failed with %v", i+1, err)
		}
		if err = ntr.Validate(); err != tc.err {
			t.Fatalf("%d: Invalid noncurrent version transition %v: err %v", i+1, ntr, err)
		}
	}
}
