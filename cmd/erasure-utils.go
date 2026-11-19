// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/klauspost/reedsolomon"
)

// getDataBlockLen - get length of data blocks from encoded blocks.
func getDataBlockLen(enBlocks [][]byte, dataBlocks int) int {
	size := 0
	// Figure out the data block length.
	for _, block := range enBlocks[:dataBlocks] {
		size += len(block)
	}
	return size
}

// Writes all the data blocks from encoded blocks until requested
// outSize length. Provides a way to skip bytes until the offset.
func writeDataBlocks(ctx context.Context, dst io.Writer, enBlocks [][]byte, dataBlocks int, offset int64, length int64) (int64, error) {
	// Offset and out size cannot be negative.
	if offset < 0 || length < 0 {
		return 0, errUnexpected
	}

	// Do we have enough blocks?
	if len(enBlocks) < dataBlocks {
		return 0, reedsolomon.ErrTooFewShards
	}

	// Do we have enough data?
	if int64(getDataBlockLen(enBlocks, dataBlocks)) < length {
		return 0, reedsolomon.ErrShortData
	}

	// Counter to decrement total left to write.
	write := length

	// Counter to increment total written.
	var totalWritten int64

	// Write all data blocks to dst.
	for _, block := range enBlocks[:dataBlocks] {
		// Skip blocks until we have reached our offset.
		if offset >= int64(len(block)) {
			// Decrement offset.
			offset -= int64(len(block))
			continue
		}

		// Skip until offset.
		block = block[offset:]

		// Reset the offset for next iteration to read everything
		// from subsequent blocks.
		offset = 0

		// We have written all the blocks, write the last remaining block.
		if write < int64(len(block)) {
			n, err := dst.Write(block[:write])
			if err != nil {
				return 0, err
			}
			totalWritten += int64(n)
			break
		}

		// Copy the block.
		n, err := dst.Write(block)
		if err != nil {
			return 0, err
		}

		// Decrement output size.
		write -= int64(n)

		// Increment written.
		totalWritten += int64(n)
	}

	// Success.
	return totalWritten, nil
}

// returns deploymentID from uploadID
func getDeplIDFromUpload(uploadID string) (string, error) {
	uploadBytes, err := base64.RawURLEncoding.DecodeString(uploadID)
	if err != nil {
		return "", fmt.Errorf("error parsing uploadID %s (%w)", uploadID, err)
	}
	slc := strings.SplitN(string(uploadBytes), ".", 2)
	if len(slc) != 2 {
		return "", fmt.Errorf("uploadID %s has incorrect format", uploadID)
	}
	return slc[0], nil
}
