// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package json

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/core-laboratories/maxio/internal/s3select/sql"
)

func TestNewPReader(t *testing.T) {
	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", file.Name()))
			if err != nil {
				t.Fatal(err)
			}
			r := NewPReader(f, &ReaderArgs{})
			var record sql.Record
			for {
				record, err = r.Read(record)
				if err != nil {
					break
				}
			}
			r.Close()
			if err != io.EOF {
				t.Fatalf("Reading failed with %s, %s", err, file.Name())
			}
		})

		t.Run(file.Name()+"-close", func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", file.Name()))
			if err != nil {
				t.Fatal(err)
			}
			r := NewPReader(f, &ReaderArgs{})
			r.Close()
			var record sql.Record
			for {
				record, err = r.Read(record)
				if err != nil {
					break
				}
			}
			if err != io.EOF {
				t.Fatalf("Reading failed with %s, %s", err, file.Name())
			}
		})
	}
}

func BenchmarkPReader(b *testing.B) {
	files, err := os.ReadDir("testdata")
	if err != nil {
		b.Fatal(err)
	}
	for _, file := range files {
		b.Run(file.Name(), func(b *testing.B) {
			f, err := os.ReadFile(filepath.Join("testdata", file.Name()))
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(f)))
			b.ReportAllocs()
			b.ResetTimer()
			var record sql.Record
			for b.Loop() {
				r := NewPReader(io.NopCloser(bytes.NewBuffer(f)), &ReaderArgs{})
				for {
					record, err = r.Read(record)
					if err != nil {
						break
					}
				}
				r.Close()
				if err != io.EOF {
					b.Fatalf("Reading failed with %s, %s", err, file.Name())
				}
			}
		})
	}
}
