// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

// Package amztime implements AWS specific time parsing and deviations
package amztime

import (
	"errors"
	"net/http"
	"time"
)

// Supported amz date formats.
var amzDateFormats = []string{
	// Do not change this order, x-amz-date format is usually in
	// iso8601Format rest are meant for relaxed handling of other
	// odd SDKs that might be out there.
	"20060102T150405Z",
	time.RFC1123,
	time.RFC1123Z,
	// Add new AMZ date formats here.
}

// ErrMalformedDate always returned for dates that cannot be parsed.
var ErrMalformedDate = errors.New("malformed date")

// Parse parses date string via supported amz date formats.
func Parse(amzDateStr string) (time.Time, error) {
	for _, dateFormat := range amzDateFormats {
		amzDate, err := time.Parse(dateFormat, amzDateStr)
		if err == nil {
			return amzDate, nil
		}
	}
	return time.Time{}, ErrMalformedDate
}

var httpTimeFormats = []string{
	// Do not change this order, http time format dates
	// are usually in http.TimeFormat however there are
	// situations where for example aws-sdk-java doesn't
	// send the correct format.
	http.TimeFormat,
	"Mon, 2 Jan 2006 15:04:05 GMT",
}

// ParseHeader parses http.TimeFormat with an acceptable
// extension for http.TimeFormat - return time might be zero
// if the timeStr is invalid.
func ParseHeader(timeStr string) (time.Time, error) {
	for _, dateFormat := range httpTimeFormats {
		t, err := time.Parse(dateFormat, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, ErrMalformedDate
}

// ParseReplicationTS parse http.TimeFormat first
// will try time.RFC3339Nano when parse http.TimeFormat failed
func ParseReplicationTS(str string) (time.Time, error) {
	tm, err := time.Parse(http.TimeFormat, str)
	if tm.IsZero() || err != nil {
		tm, err = time.Parse(time.RFC3339Nano, str)
	}
	return tm, err
}
