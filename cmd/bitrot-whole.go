// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"hash"
	"io"
)

// Implementation to calculate bitrot for the whole file.
type wholeBitrotWriter struct {
	disk      StorageAPI
	volume    string
	filePath  string
	shardSize int64 // This is the shard size of the erasure logic
	hash.Hash       // For bitrot hash
}

func (b *wholeBitrotWriter) Write(p []byte) (int, error) {
	err := b.disk.AppendFile(context.TODO(), b.volume, b.filePath, p)
	if err != nil {
		return 0, err
	}
	_, err = b.Hash.Write(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (b *wholeBitrotWriter) Close() error {
	return nil
}

// Returns whole-file bitrot writer.
func newWholeBitrotWriter(disk StorageAPI, volume, filePath string, algo BitrotAlgorithm, shardSize int64) io.WriteCloser {
	return &wholeBitrotWriter{disk, volume, filePath, shardSize, algo.New()}
}

// Implementation to verify bitrot for the whole file.
type wholeBitrotReader struct {
	disk       StorageAPI
	volume     string
	filePath   string
	verifier   *BitrotVerifier // Holds the bit-rot info
	tillOffset int64           // Affects the length of data requested in disk.ReadFile depending on Read()'s offset
	buf        []byte          // Holds bit-rot verified data
}

func (b *wholeBitrotReader) ReadAt(buf []byte, offset int64) (n int, err error) {
	if b.buf == nil {
		b.buf = make([]byte, b.tillOffset-offset)
		if _, err := b.disk.ReadFile(context.TODO(), b.volume, b.filePath, offset, b.buf, b.verifier); err != nil {
			return 0, err
		}
	}
	if len(b.buf) < len(buf) {
		return 0, errLessData
	}
	n = copy(buf, b.buf)
	b.buf = b.buf[n:]
	return n, nil
}

// Returns whole-file bitrot reader.
func newWholeBitrotReader(disk StorageAPI, volume, filePath string, algo BitrotAlgorithm, tillOffset int64, sum []byte) *wholeBitrotReader {
	return &wholeBitrotReader{
		disk:       disk,
		volume:     volume,
		filePath:   filePath,
		verifier:   &BitrotVerifier{algo, sum},
		tillOffset: tillOffset,
		buf:        nil,
	}
}
