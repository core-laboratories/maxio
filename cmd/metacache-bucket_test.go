// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"fmt"
	"testing"
)

func Benchmark_bucketMetacache_findCache(b *testing.B) {
	bm := newBucketMetacache("", false)
	const elements = 50000
	const paths = 100
	if elements%paths != 0 {
		b.Fatal("elements must be divisible by the number of paths")
	}
	var pathNames [paths]string
	for i := range pathNames[:] {
		pathNames[i] = fmt.Sprintf("prefix/%d", i)
	}
	for i := range elements {
		bm.findCache(listPathOptions{
			ID:           mustGetUUID(),
			Bucket:       "",
			BaseDir:      pathNames[i%paths],
			Prefix:       "",
			FilterPrefix: "",
			Marker:       "",
			Limit:        0,
			AskDisks:     "strict",
			Recursive:    false,
			Separator:    slashSeparator,
			Create:       true,
		})
	}
	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		bm.findCache(listPathOptions{
			ID:           mustGetUUID(),
			Bucket:       "",
			BaseDir:      pathNames[i%paths],
			Prefix:       "",
			FilterPrefix: "",
			Marker:       "",
			Limit:        0,
			AskDisks:     "strict",
			Recursive:    false,
			Separator:    slashSeparator,
			Create:       true,
		})
	}
}
