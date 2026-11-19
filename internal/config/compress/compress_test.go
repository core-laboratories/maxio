// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package compress

import (
	"reflect"
	"testing"
)

func TestParseCompressIncludes(t *testing.T) {
	testCases := []struct {
		str              string
		expectedPatterns []string
		success          bool
	}{
		// invalid input
		{",,,", []string{}, false},
		{"", []string{}, false},
		{",", []string{}, false},
		{"/", []string{}, false},
		{"text/*,/", []string{}, false},

		// valid input
		{".txt,.log", []string{".txt", ".log"}, true},
		{"text/*,application/json", []string{"text/*", "application/json"}, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.str, func(t *testing.T) {
			gotPatterns, err := parseCompressIncludes(testCase.str)
			if !testCase.success && err == nil {
				t.Error("expected failure but success instead")
			}
			if testCase.success && err != nil {
				t.Errorf("expected success but failed instead %s", err)
			}
			if testCase.success && !reflect.DeepEqual(testCase.expectedPatterns, gotPatterns) {
				t.Errorf("expected patterns %s but got %s", testCase.expectedPatterns, gotPatterns)
			}
		})
	}
}
