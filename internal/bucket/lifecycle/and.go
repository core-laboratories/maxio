// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lifecycle

import (
	"encoding/xml"
)

var errDuplicateTagKey = Errorf("Duplicate Tag Keys are not allowed")

// And - a tag to combine a prefix and multiple tags for lifecycle configuration rule.
type And struct {
	XMLName               xml.Name `xml:"And"`
	ObjectSizeGreaterThan int64    `xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    int64    `xml:"ObjectSizeLessThan,omitempty"`
	Prefix                Prefix   `xml:"Prefix,omitempty"`
	Tags                  []Tag    `xml:"Tag,omitempty"`
}

// isEmpty returns true if Tags field is null
func (a And) isEmpty() bool {
	return len(a.Tags) == 0 && !a.Prefix.set &&
		a.ObjectSizeGreaterThan == 0 && a.ObjectSizeLessThan == 0
}

// Validate - validates the And field
func (a And) Validate() error {
	// > This is used in a Lifecycle Rule Filter to apply a logical AND to two or more predicates.
	// ref: https://docs.aws.amazon.com/AmazonS3/latest/API/API_LifecycleRuleAndOperator.html
	// i.e, predCount >= 2
	var predCount int
	if a.Prefix.set {
		predCount++
	}
	predCount += len(a.Tags)
	if a.ObjectSizeGreaterThan > 0 {
		predCount++
	}
	if a.ObjectSizeLessThan > 0 {
		predCount++
	}

	if predCount < 2 {
		return errXMLNotWellFormed
	}

	if a.ContainsDuplicateTag() {
		return errDuplicateTagKey
	}
	for _, t := range a.Tags {
		if err := t.Validate(); err != nil {
			return err
		}
	}

	if a.ObjectSizeGreaterThan < 0 || a.ObjectSizeLessThan < 0 {
		return errXMLNotWellFormed
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

// BySize returns true when sz satisfies a
// ObjectSizeLessThan/ObjectSizeGreaterthan or a logical AND of these predicates
// Note: And combines size and other predicates like Tags, Prefix, etc. This
// method applies exclusively to size predicates only.
func (a And) BySize(sz int64) bool {
	if a.ObjectSizeGreaterThan > 0 &&
		sz <= a.ObjectSizeGreaterThan {
		return false
	}
	if a.ObjectSizeLessThan > 0 &&
		sz >= a.ObjectSizeLessThan {
		return false
	}
	return true
}
