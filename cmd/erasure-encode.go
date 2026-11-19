// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"fmt"
	"io"
)

// Writes to multiple writers
type multiWriter struct {
	writers     []io.Writer
	writeQuorum int
	errs        []error
}

// Write writes data to writers.
func (p *multiWriter) Write(ctx context.Context, blocks [][]byte) error {
	for i := range p.writers {
		if p.errs[i] != nil {
			continue
		}
		if p.writers[i] == nil {
			p.errs[i] = errDiskNotFound
			continue
		}
		var n int
		n, p.errs[i] = p.writers[i].Write(blocks[i])
		if p.errs[i] == nil {
			if n != len(blocks[i]) {
				p.errs[i] = io.ErrShortWrite
				p.writers[i] = nil
			}
		} else {
			p.writers[i] = nil
		}
	}

	// If nilCount >= p.writeQuorum, we return nil. This is because HealFile() uses
	// CreateFile with p.writeQuorum=1 to accommodate healing of single disk.
	// i.e if we do no return here in such a case, reduceWriteQuorumErrs() would
	// return a quorum error to HealFile().
	nilCount := countErrs(p.errs, nil)
	if nilCount >= p.writeQuorum {
		return nil
	}

	writeErr := reduceWriteQuorumErrs(ctx, p.errs, objectOpIgnoredErrs, p.writeQuorum)
	return fmt.Errorf("%w (offline-disks=%d/%d)", writeErr, countErrs(p.errs, errDiskNotFound), len(p.writers))
}

// Encode reads from the reader, erasure-encodes the data and writes to the writers.
func (e *Erasure) Encode(ctx context.Context, src io.Reader, writers []io.Writer, buf []byte, quorum int) (total int64, err error) {
	writer := &multiWriter{
		writers:     writers,
		writeQuorum: quorum,
		errs:        make([]error, len(writers)),
	}

	for {
		var blocks [][]byte
		n, err := io.ReadFull(src, buf)
		if err != nil {
			if !IsErrIgnored(err, []error{
				io.EOF,
				io.ErrUnexpectedEOF,
			}...) {
				return 0, err
			}
		}

		eof := err == io.EOF || err == io.ErrUnexpectedEOF
		if n == 0 && total != 0 {
			// Reached EOF, nothing more to be done.
			break
		}

		// We take care of the situation where if n == 0 and total == 0 by creating empty data and parity files.
		blocks, err = e.EncodeData(ctx, buf[:n])
		if err != nil {
			return 0, err
		}

		if err = writer.Write(ctx, blocks); err != nil {
			return 0, err
		}

		total += int64(n)
		if eof {
			break
		}
	}
	return total, nil
}
