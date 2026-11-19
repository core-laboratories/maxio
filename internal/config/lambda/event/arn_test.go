// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import (
	"testing"
)

func TestARNString(t *testing.T) {
	testCases := []struct {
		arn            ARN
		expectedResult string
	}{
		{ARN{}, ""},
		{ARN{TargetID{"1", "webhook"}, ""}, "arn:minio:s3-object-lambda::1:webhook"},
		{ARN{TargetID{"1", "webhook"}, "us-east-1"}, "arn:minio:s3-object-lambda:us-east-1:1:webhook"},
	}

	for i, testCase := range testCases {
		result := testCase.arn.String()

		if result != testCase.expectedResult {
			t.Fatalf("test %v: result: expected: %v, got: %v", i+1, testCase.expectedResult, result)
		}
	}
}

func TestParseARN(t *testing.T) {
	testCases := []struct {
		s           string
		expectedARN *ARN
		expectErr   bool
	}{
		{"", nil, true},
		{"arn:minio:s3-object-lambda:::", nil, true},
		{"arn:minio:s3-object-lambda::1:webhook:remote", nil, true},
		{"arn:aws:s3-object-lambda::1:webhook", nil, true},
		{"arn:minio:sns::1:webhook", nil, true},
		{"arn:minio:s3-object-lambda::1:webhook", &ARN{TargetID{"1", "webhook"}, ""}, false},
		{"arn:minio:s3-object-lambda:us-east-1:1:webhook", &ARN{TargetID{"1", "webhook"}, "us-east-1"}, false},
	}

	for i, testCase := range testCases {
		arn, err := ParseARN(testCase.s)
		expectErr := (err != nil)

		if expectErr != testCase.expectErr {
			t.Fatalf("test %v: error: expected: %v, got: %v", i+1, testCase.expectErr, expectErr)
		}

		if !testCase.expectErr {
			if *arn != *testCase.expectedARN {
				t.Fatalf("test %v: data: expected: %v, got: %v", i+1, testCase.expectedARN, arn)
			}
		}
	}
}
