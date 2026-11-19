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
	"testing"

	"github.com/minio/mux"
)

// Test cross domain xml handler.
func TestCrossXMLHandler(t *testing.T) {
	// Server initialization.
	router := mux.NewRouter().SkipClean(true).UseEncodedPath()
	handler := setCrossDomainPolicyMiddleware(router)
	srv := httptest.NewServer(handler)

	resp, err := http.Get(srv.URL + crossDomainXMLEntity)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected http status received", resp.Status)
	}
}
