// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"io"
	"testing"
)

func testBitrotReaderWriterAlgo(t *testing.T, bitrotAlgo BitrotAlgorithm) {
	tmpDir := t.TempDir()

	volume := "testvol"
	filePath := "testfile"

	disk, err := newLocalXLStorage(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	disk.MakeVol(t.Context(), volume)

	writer := newBitrotWriter(disk, "", volume, filePath, 35, bitrotAlgo, 10)

	_, err = writer.Write([]byte("aaaaaaaaaa"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = writer.Write([]byte("aaaaaaaaaa"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = writer.Write([]byte("aaaaaaaaaa"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = writer.Write([]byte("aaaaa"))
	if err != nil {
		t.Fatal(err)
	}
	if bw, ok := writer.(io.Closer); ok {
		bw.Close()
	}

	reader := newBitrotReader(disk, nil, volume, filePath, 35, bitrotAlgo, bitrotWriterSum(writer), 10)
	b := make([]byte, 10)
	if _, err = reader.ReadAt(b, 0); err != nil {
		t.Fatal(err)
	}
	if _, err = reader.ReadAt(b, 10); err != nil {
		t.Fatal(err)
	}
	if _, err = reader.ReadAt(b, 20); err != nil {
		t.Fatal(err)
	}
	if _, err = reader.ReadAt(b[:5], 30); err != nil {
		t.Fatal(err)
	}
	if br, ok := reader.(io.Closer); ok {
		br.Close()
	}
}

func TestAllBitrotAlgorithms(t *testing.T) {
	for bitrotAlgo := range bitrotAlgorithms {
		testBitrotReaderWriterAlgo(t, bitrotAlgo)
	}
}
