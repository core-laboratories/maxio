// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func Test_readFromSecret(t *testing.T) {
	testCases := []struct {
		content       string
		expectedErr   bool
		expectedValue string
	}{
		{
			"value\n",
			false,
			"value",
		},
		{
			" \t\n Hello, Gophers \n\t\r\n",
			false,
			"Hello, Gophers",
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			tmpfile, err := os.CreateTemp(t.TempDir(), "testfile")
			if err != nil {
				t.Error(err)
			}
			tmpfile.WriteString(testCase.content)
			tmpfile.Sync()
			tmpfile.Close()

			value, err := readFromSecret(tmpfile.Name())
			if err != nil && !testCase.expectedErr {
				t.Error(err)
			}
			if err == nil && testCase.expectedErr {
				t.Error(errors.New("expected error, found success"))
			}
			if value != testCase.expectedValue {
				t.Errorf("Expected %s, got %s", testCase.expectedValue, value)
			}
		})
	}
}

func Test_minioEnvironFromFile(t *testing.T) {
	testCases := []struct {
		content      string
		expectedErr  bool
		expectedEkvs []envKV
	}{
		{
			`
export MINIO_ROOT_USER=minio
export MINIO_ROOT_PASSWORD=minio123`,
			false,
			[]envKV{
				{
					Key:   "MINIO_ROOT_USER",
					Value: "minio",
				},
				{
					Key:   "MINIO_ROOT_PASSWORD",
					Value: "minio123",
				},
			},
		},
		// Value with double quotes
		{
			`export MINIO_ROOT_USER="minio"`,
			false,
			[]envKV{
				{
					Key:   "MINIO_ROOT_USER",
					Value: "minio",
				},
			},
		},
		// Value with single quotes
		{
			`export MINIO_ROOT_USER='minio'`,
			false,
			[]envKV{
				{
					Key:   "MINIO_ROOT_USER",
					Value: "minio",
				},
			},
		},
		{
			`
MINIO_ROOT_USER=minio
MINIO_ROOT_PASSWORD=minio123`,
			false,
			[]envKV{
				{
					Key:   "MINIO_ROOT_USER",
					Value: "minio",
				},
				{
					Key:   "MINIO_ROOT_PASSWORD",
					Value: "minio123",
				},
			},
		},
		{
			`
export MINIO_ROOT_USERminio
export MINIO_ROOT_PASSWORD=minio123`,
			true,
			nil,
		},
		{
			`
# simple comment
# MINIO_ROOT_USER=minioadmin
# MINIO_ROOT_PASSWORD=minioadmin
MINIO_ROOT_USER=minio
MINIO_ROOT_PASSWORD=minio123`,
			false,
			[]envKV{
				{
					Key:   "MINIO_ROOT_USER",
					Value: "minio",
				},
				{
					Key:   "MINIO_ROOT_PASSWORD",
					Value: "minio123",
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			tmpfile, err := os.CreateTemp(t.TempDir(), "testfile")
			if err != nil {
				t.Error(err)
			}
			tmpfile.WriteString(testCase.content)
			tmpfile.Sync()
			tmpfile.Close()

			ekvs, err := minioEnvironFromFile(tmpfile.Name())
			if err != nil && !testCase.expectedErr {
				t.Error(err)
			}
			if err == nil && testCase.expectedErr {
				t.Error(errors.New("expected error, found success"))
			}

			if len(ekvs) != len(testCase.expectedEkvs) {
				t.Errorf("expected %v keys, got %v keys", len(testCase.expectedEkvs), len(ekvs))
			}

			if !reflect.DeepEqual(ekvs, testCase.expectedEkvs) {
				t.Errorf("expected %v, got %v", testCase.expectedEkvs, ekvs)
			}
		})
	}
}
