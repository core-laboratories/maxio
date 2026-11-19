// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package replication

import (
	"encoding/xml"
)

// And - a tag to combine a prefix and multiple tags for replication configuration rule.
type And struct {
	XMLName xml.Name `xml:"And" json:"And"`
	Prefix  string   `xml:"Prefix,omitempty" json:"Prefix,omitempty"`
	Tags    []Tag    `xml:"Tag,omitempty" json:"Tag,omitempty"`
}

var errDuplicateTagKey = Errorf("Duplicate Tag Keys are not allowed")

// isEmpty returns true if Tags field is null
func (a And) isEmpty() bool {
	return len(a.Tags) == 0 && a.Prefix == ""
}

// Validate - validates the And field
func (a And) Validate() error {
	if a.ContainsDuplicateTag() {
		return errDuplicateTagKey
	}
	for _, t := range a.Tags {
		if err := t.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ContainsDuplicateTag - returns true if duplicate keys are present in And
func (a And) ContainsDuplicateTag() bool {
	x := make(map[string]struct{}, len(a.Tags))

	for _, t := range a.Tags {
		if _, has := x[t.Key]; has {
			return true
		}
		x[t.Key] = struct{}{}
	}

	return false
}
