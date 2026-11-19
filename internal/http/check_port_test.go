// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package http

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"testing"
)

// Tests for port availability logic written for server startup sequence.
func TestCheckPortAvailability(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip()
	}

	l, err := net.Listen("tcp", "localhost:0") // ask kernel for a free port.
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	testCases := []struct {
		host        string
		port        int
		expectedErr error
	}{
		{"", port, fmt.Errorf("listen tcp :%v: bind: address already in use", port)},
		{"127.0.0.1", port, fmt.Errorf("listen tcp 127.0.0.1:%v: bind: address already in use", port)},
	}

	for _, testCase := range testCases {
		err := CheckPortAvailability(testCase.host, strconv.Itoa(testCase.port), TCPOptions{})
		switch {
		case testCase.expectedErr == nil:
			if err != nil {
				t.Fatalf("error: expected = <nil>, got = %v", err)
			}
		case err == nil:
			t.Fatalf("error: expected = %v, got = <nil>", testCase.expectedErr)
		case testCase.expectedErr.Error() != err.Error():
			t.Fatalf("error: expected = %v, got = %v", testCase.expectedErr, err)
		}
	}
}
