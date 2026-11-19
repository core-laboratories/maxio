// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package store

import (
	"testing"
)

func TestKeyString(t *testing.T) {
	testCases := []struct {
		key            Key
		expectedString string
	}{
		{
			key: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Extension: ".event",
			},
			expectedString: "01894394-d046-4783-ba0d-f1c6885790dc.event",
		},
		{
			key: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Compress:  true,
				Extension: ".event",
				ItemCount: 100,
			},
			expectedString: "100:01894394-d046-4783-ba0d-f1c6885790dc.event.snappy",
		},
		{
			key: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Extension: ".event",
				ItemCount: 100,
			},
			expectedString: "100:01894394-d046-4783-ba0d-f1c6885790dc.event",
		},
		{
			key: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Compress:  true,
				Extension: ".event",
				ItemCount: 1,
			},
			expectedString: "01894394-d046-4783-ba0d-f1c6885790dc.event.snappy",
		},
		{
			key: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Extension: ".event",
				ItemCount: 1,
			},
			expectedString: "01894394-d046-4783-ba0d-f1c6885790dc.event",
		},
	}

	for i, testCase := range testCases {
		if testCase.expectedString != testCase.key.String() {
			t.Fatalf("case[%v]: key.String() Expected: %s, got %s", i, testCase.expectedString, testCase.key.String())
		}
	}
}

func TestParseKey(t *testing.T) {
	testCases := []struct {
		str         string
		expectedKey Key
	}{
		{
			str: "01894394-d046-4783-ba0d-f1c6885790dc.event",
			expectedKey: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Extension: ".event",
				ItemCount: 1,
			},
		},
		{
			str: "100:01894394-d046-4783-ba0d-f1c6885790dc.event.snappy",
			expectedKey: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Compress:  true,
				Extension: ".event",
				ItemCount: 100,
			},
		},
		{
			str: "100:01894394-d046-4783-ba0d-f1c6885790dc.event",
			expectedKey: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Extension: ".event",
				ItemCount: 100,
			},
		},
		{
			str: "01894394-d046-4783-ba0d-f1c6885790dc.event.snappy",
			expectedKey: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Compress:  true,
				Extension: ".event",
				ItemCount: 1,
			},
		},
		{
			str: "01894394-d046-4783-ba0d-f1c6885790dc.event",
			expectedKey: Key{
				Name:      "01894394-d046-4783-ba0d-f1c6885790dc",
				Extension: ".event",
				ItemCount: 1,
			},
		},
	}

	for i, testCase := range testCases {
		key := parseKey(testCase.str)
		if testCase.expectedKey.Name != key.Name {
			t.Fatalf("case[%v]: Expected key.Name: %v, got %v", i, testCase.expectedKey.Name, key.Name)
		}
		if testCase.expectedKey.Compress != key.Compress {
			t.Fatalf("case[%v]: Expected key.Compress: %v, got %v", i, testCase.expectedKey.Compress, key.Compress)
		}
		if testCase.expectedKey.Extension != key.Extension {
			t.Fatalf("case[%v]: Expected key.Extension: %v, got %v", i, testCase.expectedKey.Extension, key.Extension)
		}
		if testCase.expectedKey.ItemCount != key.ItemCount {
			t.Fatalf("case[%v]: Expected key.ItemCount: %v, got %v", i, testCase.expectedKey.ItemCount, key.ItemCount)
		}
		if testCase.expectedKey.String() != key.String() {
			t.Fatalf("case[%v]: Expected key.String(): %v, got %v", i, testCase.expectedKey.String(), key.String())
		}
	}
}
