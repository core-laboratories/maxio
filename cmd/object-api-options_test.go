// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	xhttp "github.com/core-laboratories/maxio/internal/http"
)

// TestGetAndValidateAttributesOpts is currently minimal and covers a subset of getAndValidateAttributesOpts(),
// it is intended to be expanded when the function is worked on in the future.
func TestGetAndValidateAttributesOpts(t *testing.T) {
	globalBucketVersioningSys = &BucketVersioningSys{}
	bucket := minioMetaBucket
	ctx := t.Context()
	testCases := []struct {
		name            string
		headers         http.Header
		wantObjectAttrs map[string]struct{}
	}{
		{
			name:            "empty header",
			headers:         http.Header{},
			wantObjectAttrs: map[string]struct{}{},
		},
		{
			name: "single header line",
			headers: http.Header{
				xhttp.AmzObjectAttributes: []string{"test1,test2"},
			},
			wantObjectAttrs: map[string]struct{}{
				"test1": {}, "test2": {},
			},
		},
		{
			name: "multiple header lines with some duplicates",
			headers: http.Header{
				xhttp.AmzObjectAttributes: []string{"test1,test2", "test3,test4", "test4,test3"},
			},
			wantObjectAttrs: map[string]struct{}{
				"test1": {}, "test2": {}, "test3": {}, "test4": {},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header = testCase.headers

			opts, _ := getAndValidateAttributesOpts(ctx, rec, req, bucket, "testobject")

			if !reflect.DeepEqual(opts.ObjectAttributes, testCase.wantObjectAttrs) {
				t.Errorf("want opts %v, got %v", testCase.wantObjectAttrs, opts.ObjectAttributes)
			}
		})
	}
}
