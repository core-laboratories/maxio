// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package amztime

import (
	"testing"
	"time"
)

func TestISO8601Format(t *testing.T) {
	testCases := []struct {
		date           time.Time
		expectedOutput string
	}{
		{
			date:           time.Date(2009, time.November, 13, 4, 51, 1, 940303531, time.UTC),
			expectedOutput: "2009-11-13T04:51:01.940Z",
		},
		{
			date:           time.Date(2009, time.November, 13, 4, 51, 1, 901303531, time.UTC),
			expectedOutput: "2009-11-13T04:51:01.901Z",
		},
		{
			date:           time.Date(2009, time.November, 13, 4, 51, 1, 900303531, time.UTC),
			expectedOutput: "2009-11-13T04:51:01.900Z",
		},
		{
			date:           time.Date(2009, time.November, 13, 4, 51, 1, 941303531, time.UTC),
			expectedOutput: "2009-11-13T04:51:01.941Z",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedOutput, func(t *testing.T) {
			gotOutput := ISO8601Format(testCase.date)
			t.Log("Go", testCase.date.Format(iso8601TimeFormat))
			if gotOutput != testCase.expectedOutput {
				t.Errorf("Expected %s, got %s", testCase.expectedOutput, gotOutput)
			}
		})
	}
}
