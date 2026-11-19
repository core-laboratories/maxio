// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/minio/madmin-go/v3/estream"
)

type keepFileErr struct {
	error
}

func extractInspectV2(pks [][]byte, r io.Reader, extractDir string) error {
	sr, err := estream.NewReader(r)
	if err != nil {
		return err
	}
	for _, pk := range pks {
		privKey, err := bytesToPrivateKey(pk)
		if err != nil {
			return fmt.Errorf("decoding key returned: %w", err)
		}
		sr.SetPrivateKey(privKey)
	}
	sr.ReturnNonDecryptable(true)

	// Debug corrupted streams.
	if false {
		sr.SkipEncrypted(true)
		return sr.DebugStream(os.Stdout)
	}
	extracted := false
	for {
		stream, err := sr.NextStream()
		if err != nil {
			if err == io.EOF {
				if extracted {
					return nil
				}
				return errors.New("no data found on stream")
			}
			if errors.Is(err, estream.ErrNoKey) {
				fmt.Println("Skipping", stream.Name, "no private key")
				if err := stream.Skip(); err != nil {
					return fmt.Errorf("stream skip: %w", err)
				}
				continue
			}
			return fmt.Errorf("next stream: %w", err)
		}
		if strings.Contains(stream.Name, "..") || !utf8.ValidString(stream.Name) {
			return fmt.Errorf("invalid stream name: %q", stream.Name)
		}

		dst := filepath.Join(extractDir, stream.Name)
		os.Mkdir(extractDir, 0o755)
		w, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		_, err = io.Copy(w, stream)
		if err != nil {
			return fmt.Errorf("reading inspect stream: %w", err)
		}
		fmt.Printf("Extracted: %s\n", dst)
		extracted = true
	}
}
