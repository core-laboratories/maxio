// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"fmt"
	"testing"
)

func TestS3EncodeName(t *testing.T) {
	testCases := []struct {
		inputText, encodingType, expectedOutput string
	}{
		{"a b", "", "a b"},
		{"a b", "url", "a+b"},
		{"p- ", "url", "p-+"},
		{"p-%", "url", "p-%25"},
		{"p/", "url", "p/"},
		{"p/", "url", "p/"},
		{"~user", "url", "%7Euser"},
		{"*user", "url", "*user"},
		{"user+password", "url", "user%2Bpassword"},
		{"_user", "url", "_user"},
		{"firstname.lastname", "url", "firstname.lastname"},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			outputText := s3EncodeName(testCase.inputText, testCase.encodingType)
			if testCase.expectedOutput != outputText {
				t.Errorf("Expected `%s`, got `%s`", testCase.expectedOutput, outputText)
			}
		})
	}
}
