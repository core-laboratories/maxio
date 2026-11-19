// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package kms

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestSingleKeyRoundtrip(t *testing.T) {
	KMS, err := ParseSecretKey("my-key:eEm+JI9/q4JhH8QwKvf3LKo4DEBl6QbfvAl1CAbMIv8=")
	if err != nil {
		t.Fatalf("Failed to initialize KMS: %v", err)
	}

	key, err := KMS.GenerateKey(t.Context(), &GenerateKeyRequest{Name: "my-key"})
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	plaintext, err := KMS.Decrypt(t.Context(), &DecryptRequest{
		Name:       key.KeyID,
		Ciphertext: key.Ciphertext,
	})
	if err != nil {
		t.Fatalf("Failed to decrypt key: %v", err)
	}
	if !bytes.Equal(key.Plaintext, plaintext) {
		t.Fatalf("Decrypted key does not match generated one: got %x - want %x", key.Plaintext, plaintext)
	}
}

func TestDecryptKey(t *testing.T) {
	KMS, err := ParseSecretKey("my-key:eEm+JI9/q4JhH8QwKvf3LKo4DEBl6QbfvAl1CAbMIv8=")
	if err != nil {
		t.Fatalf("Failed to initialize KMS: %v", err)
	}

	for i, test := range decryptKeyTests {
		dataKey, err := base64.StdEncoding.DecodeString(test.Plaintext)
		if err != nil {
			t.Fatalf("Test %d: failed to decode plaintext key: %v", i, err)
		}
		plaintext, err := KMS.Decrypt(t.Context(), &DecryptRequest{
			Name:           test.KeyID,
			Ciphertext:     []byte(test.Ciphertext),
			AssociatedData: test.Context,
		})
		if err != nil {
			t.Fatalf("Test %d: failed to decrypt key: %v", i, err)
		}
		if !bytes.Equal(plaintext, dataKey) {
			t.Fatalf("Test %d: decrypted key does not generated one: got %x - want %x", i, plaintext, dataKey)
		}
	}
}

var decryptKeyTests = []struct {
	KeyID      string
	Plaintext  string
	Ciphertext string
	Context    Context
}{
	{
		KeyID:      "my-key",
		Plaintext:  "zmS7NrG765UZ0ZN85oPjybelxqVvpz01vxsSpOISy2M=",
		Ciphertext: `{"aead":"ChaCha20Poly1305","iv":"JbI+vwvYww1lCb5VpkAFuQ==","nonce":"ARjIjJxBSD541Gz8","bytes":"KCbEc2sA0TLvA7aWTWa23AdccVfJMpOxwgG8hm+4PaNrxYfy1xFWZg2gEenVrOgv"}`,
	},
	{
		KeyID:      "my-key",
		Plaintext:  "UnPWsZgVI+T4L9WGNzFlP1PsP1Z6hn2Fx8ISeZfDGnA=",
		Ciphertext: `{"aead":"ChaCha20Poly1305","iv":"r4+yfiVbVIYR0Z2I9Fq+6g==","nonce":"2YpwGwE59GcVraI3","bytes":"k/svMglOU7/Kgwv73heG38NWW575XLcFp3SaxQHDMjJGYyRI3Fiygu2OeutGPXNL"}`,
		Context:    Context{"key": "value"},
	},
}
