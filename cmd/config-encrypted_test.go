// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"bytes"
	"testing"

	"github.com/minio/madmin-go/v3"
	"github.com/core-laboratories/maxio/internal/auth"
)

func TestDecryptData(t *testing.T) {
	cred1 := auth.Credentials{
		AccessKey: "minio",
		SecretKey: "minio123",
	}

	cred2 := auth.Credentials{
		AccessKey: "minio",
		SecretKey: "minio1234",
	}

	data := []byte(`config data`)
	edata1, err := madmin.EncryptData(cred1.String(), data)
	if err != nil {
		t.Fatal(err)
	}

	edata2, err := madmin.EncryptData(cred2.String(), data)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		edata   []byte
		cred    auth.Credentials
		success bool
	}{
		{edata1, cred1, true},
		{edata2, cred2, true},
		{data, cred1, false},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			ddata, err := madmin.DecryptData(test.cred.String(), bytes.NewReader(test.edata))
			if err != nil && test.success {
				t.Errorf("Expected success, saw failure %v", err)
			}
			if err == nil && !test.success {
				t.Error("Expected failure, saw success")
			}
			if test.success {
				if !bytes.Equal(ddata, data) {
					t.Errorf("Expected %s, got %s", string(data), string(ddata))
				}
			}
		})
	}
}
