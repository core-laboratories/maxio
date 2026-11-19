// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// markerTagVersion is the marker version.
// Should not need to be updated unless a fundamental change is made to the marker format.
const markerTagVersion = "v2"

// parseMarker will parse a marker possibly encoded with encodeMarker
func (o *listPathOptions) parseMarker() {
	s := o.Marker
	if !strings.Contains(s, "[minio_cache:"+markerTagVersion) {
		return
	}
	start := strings.LastIndex(s, "[")
	o.Marker = s[:start]
	end := strings.LastIndex(s, "]")
	tag := strings.Trim(s[start:end], "[]")
	tags := strings.SplitSeq(tag, ",")
	for tag := range tags {
		kv := strings.Split(tag, ":")
		if len(kv) < 2 {
			continue
		}
		switch kv[0] {
		case "minio_cache":
			if kv[1] != markerTagVersion {
				continue
			}
		case "id":
			o.ID = kv[1]
		case "return":
			o.ID = mustGetUUID()
			o.Create = true
		case "p": // pool
			v, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				o.ID = mustGetUUID()
				o.Create = true
				continue
			}
			o.pool = int(v)
		case "s": // set
			v, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				o.ID = mustGetUUID()
				o.Create = true
				continue
			}
			o.set = int(v)
		default:
			// Ignore unknown
		}
	}
}

// encodeMarker will encode a uuid and return it as a marker.
// uuid cannot contain '[', ':' or ','.
func (o listPathOptions) encodeMarker(marker string) string {
	if o.ID == "" {
		// Mark as returning listing...
		return fmt.Sprintf("%s[minio_cache:%s,return:]", marker, markerTagVersion)
	}
	if strings.ContainsAny(o.ID, "[:,") {
		internalLogIf(context.Background(), fmt.Errorf("encodeMarker: uuid %s contained invalid characters", o.ID))
	}
	return fmt.Sprintf("%s[minio_cache:%s,id:%s,p:%d,s:%d]", marker, markerTagVersion, o.ID, o.pool, o.set)
}
