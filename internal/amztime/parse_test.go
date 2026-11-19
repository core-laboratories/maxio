// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

// Package amztime implements AWS specific time parsing and deviations
package amztime

import (
	"errors"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type testCase struct {
		expectedErr  error
		expectedTime time.Time
		timeStr      string
	}
	testCases := []testCase{
		{
			ErrMalformedDate,
			time.Time{},
			"Tue Sep  6 07:10:23 PM PDT 2022",
		},
		{
			nil,
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			"Tue, 10 Nov 2009 23:00:00 UTC",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.timeStr, func(t *testing.T) {
			gott, goterr := Parse(testCase.timeStr)
			if !errors.Is(goterr, testCase.expectedErr) {
				t.Errorf("expected %v, got %v", testCase.expectedErr, goterr)
			}
			if !gott.Equal(testCase.expectedTime) {
				t.Errorf("expected %v, got %v", testCase.expectedTime, gott)
			}
		})
	}
}
