// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/secure-io/sio-go"
)

func extractInspectV1(keyHex string, r io.Reader, w io.Writer, okMsg string) error {
	id, err := hex.DecodeString(keyHex[:8])
	if err != nil {
		return err
	}
	key, err := hex.DecodeString(keyHex[8:])
	if err != nil {
		return err
	}
	// Verify that CRC is ok.
	want := binary.LittleEndian.Uint32(id)
	got := crc32.ChecksumIEEE(key)
	if want != got {
		return fmt.Errorf("Invalid key checksum, want %x, got %x", want, got)
	}

	stream, err := sio.AES_256_GCM.Stream(key)
	if err != nil {
		return err
	}
	// Zero nonce, we only use each key once, and 32 bytes is plenty.
	nonce := make([]byte, stream.NonceSize())
	encr := stream.DecryptReader(r, nonce, nil)
	_, err = io.Copy(w, encr)
	if err == nil {
		fmt.Println(okMsg)
	}
	return err
}
