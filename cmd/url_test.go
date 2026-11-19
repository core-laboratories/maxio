// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"net/http"
	"testing"
)

func BenchmarkURLQueryForm(b *testing.B) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:9000/bucket/name?uploadId=upload&partNumber=1", http.NoBody)
	if err != nil {
		b.Fatal(err)
	}

	// benchmark utility which helps obtain number of allocations and bytes allocated per ops.
	b.ReportAllocs()
	// the actual benchmark for PutObject starts here. Reset the benchmark timer.
	b.ResetTimer()

	if err := req.ParseForm(); err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req.Form.Get("uploadId")
		}
	})

	// Benchmark ends here. Stop timer.
	b.StopTimer()
}

// BenchmarkURLQuery - benchmark URL memory allocations
func BenchmarkURLQuery(b *testing.B) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:9000/bucket/name?uploadId=upload&partNumber=1", http.NoBody)
	if err != nil {
		b.Fatal(err)
	}

	// benchmark utility which helps obtain number of allocations and bytes allocated per ops.
	b.ReportAllocs()
	// the actual benchmark for PutObject starts here. Reset the benchmark timer.
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req.URL.Query().Get("uploadId")
		}
	})

	// Benchmark ends here. Stop timer.
	b.StopTimer()
}
