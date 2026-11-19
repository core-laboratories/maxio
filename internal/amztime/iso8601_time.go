// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package amztime

import (
	"strings"
	"time"
)

// RFC3339 a subset of the ISO8601 timestamp format. e.g 2014-04-29T18:30:38Z
const (
	iso8601TimeFormat     = "2006-01-02T15:04:05.000Z"    // Reply date format with millisecond precision.
	iso8601TimeFormatLong = "2006-01-02T15:04:05.000000Z" // Reply date format with nanosecond precision.
)

// ISO8601Format converts time 't' into ISO8601 time format expected in AWS S3 spec.
//
// This function is needed to avoid a Go's float64 precision bug, where Go avoids
// padding the extra '0' before the timezone.
func ISO8601Format(t time.Time) string {
	value := t.Format(iso8601TimeFormat)
	if len(value) < len(iso8601TimeFormat) {
		value = t.Format(iso8601TimeFormat[:len(iso8601TimeFormat)-1])
		// Pad necessary zeroes to full-fill the iso8601TimeFormat
		return value + strings.Repeat("0", (len(iso8601TimeFormat)-1)-len(value)) + "Z"
	}
	return value
}

// ISO8601Parse parses ISO8601 date string
func ISO8601Parse(iso8601 string) (t time.Time, err error) {
	for _, layout := range []string{
		iso8601TimeFormat,
		iso8601TimeFormatLong,
		time.RFC3339,
	} {
		t, err = time.Parse(layout, iso8601)
		if err == nil {
			return t, nil
		}
	}

	return t, err
}
