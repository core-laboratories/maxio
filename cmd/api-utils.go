// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '/', '*':
		return false
	}
	return true
}

// s3URLEncode is based on Golang's url.QueryEscape() code,
// while considering some S3 exceptions:
//   - Avoid encoding '/' and '*'
//   - Force encoding of '~'
func s3URLEncode(s string) string {
	spaceCount, hexCount := 0, 0
	for i := range len(s) {
		c := s[i]
		if shouldEscape(c) {
			if c == ' ' {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := range len(s) {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := range len(s) {
		switch c := s[i]; {
		case c == ' ':
			t[j] = '+'
			j++
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// s3EncodeName encodes string in response when encodingType is specified in AWS S3 requests.
func s3EncodeName(name, encodingType string) string {
	if strings.ToLower(encodingType) == "url" {
		return s3URLEncode(name)
	}
	return name
}

// getHandlerName returns the name of the handler function. It takes the type
// name as a string to clean up the name retrieved via reflection. This function
// only works correctly when the type is present in the cmd package.
func getHandlerName(f http.HandlerFunc, cmdType string) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	packageName := fmt.Sprintf("github.com/core-laboratories/maxio/cmd.%s.", cmdType)
	name = strings.TrimPrefix(name, packageName)
	name = strings.TrimSuffix(name, "Handler-fm")
	name = strings.TrimSuffix(name, "-fm")
	return name
}
