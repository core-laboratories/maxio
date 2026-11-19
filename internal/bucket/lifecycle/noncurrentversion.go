// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lifecycle

import (
	"encoding/xml"
	"time"
)

// NoncurrentVersionExpiration - an action for lifecycle configuration rule.
type NoncurrentVersionExpiration struct {
	XMLName                 xml.Name       `xml:"NoncurrentVersionExpiration"`
	NoncurrentDays          ExpirationDays `xml:"NoncurrentDays,omitempty"`
	NewerNoncurrentVersions int            `xml:"NewerNoncurrentVersions,omitempty"`
	set                     bool
}

// MarshalXML if non-current days not set to non zero value
func (n NoncurrentVersionExpiration) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.IsNull() {
		return nil
	}
	type noncurrentVersionExpirationWrapper NoncurrentVersionExpiration
	return e.EncodeElement(noncurrentVersionExpirationWrapper(n), start)
}

// UnmarshalXML decodes NoncurrentVersionExpiration
func (n *NoncurrentVersionExpiration) UnmarshalXML(d *xml.Decoder, startElement xml.StartElement) error {
	// To handle xml with MaxNoncurrentVersions from older MinIO releases.
	// note: only one of MaxNoncurrentVersions or NewerNoncurrentVersions would be present.
	type noncurrentExpiration struct {
		XMLName                 xml.Name       `xml:"NoncurrentVersionExpiration"`
		NoncurrentDays          ExpirationDays `xml:"NoncurrentDays,omitempty"`
		NewerNoncurrentVersions int            `xml:"NewerNoncurrentVersions,omitempty"`
		MaxNoncurrentVersions   int            `xml:"MaxNoncurrentVersions,omitempty"`
	}

	var val noncurrentExpiration
	err := d.DecodeElement(&val, &startElement)
	if err != nil {
		return err
	}
	if val.MaxNoncurrentVersions > 0 {
		val.NewerNoncurrentVersions = val.MaxNoncurrentVersions
	}
	*n = NoncurrentVersionExpiration{
		XMLName:                 val.XMLName,
		NoncurrentDays:          val.NoncurrentDays,
		NewerNoncurrentVersions: val.NewerNoncurrentVersions,
	}
	n.set = true
	return nil
}

// IsNull returns if both NoncurrentDays and NoncurrentVersions are empty
func (n NoncurrentVersionExpiration) IsNull() bool {
	return n.IsDaysNull() && n.NewerNoncurrentVersions == 0
}

// IsDaysNull returns true if days field is null
func (n NoncurrentVersionExpiration) IsDaysNull() bool {
	return n.NoncurrentDays == ExpirationDays(0)
}

// Validate returns an error with wrong value
func (n NoncurrentVersionExpiration) Validate() error {
	if !n.set {
		return nil
	}
	val := int(n.NoncurrentDays)
	switch {
	case val == 0 && n.NewerNoncurrentVersions == 0:
		// both fields can't be zero
		return errXMLNotWellFormed

	case val < 0, n.NewerNoncurrentVersions < 0:
		// negative values are not supported
		return errXMLNotWellFormed
	}
	return nil
}

// NoncurrentVersionTransition - an action for lifecycle configuration rule.
type NoncurrentVersionTransition struct {
	NoncurrentDays TransitionDays `xml:"NoncurrentDays"`
	StorageClass   string         `xml:"StorageClass"`
	set            bool
}

// MarshalXML is extended to leave out
// <NoncurrentVersionTransition></NoncurrentVersionTransition> tags
func (n NoncurrentVersionTransition) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.IsNull() {
		return nil
	}
	type noncurrentVersionTransitionWrapper NoncurrentVersionTransition
	return e.EncodeElement(noncurrentVersionTransitionWrapper(n), start)
}

// UnmarshalXML decodes NoncurrentVersionExpiration
func (n *NoncurrentVersionTransition) UnmarshalXML(d *xml.Decoder, startElement xml.StartElement) error {
	type noncurrentVersionTransitionWrapper NoncurrentVersionTransition
	var val noncurrentVersionTransitionWrapper
	err := d.DecodeElement(&val, &startElement)
	if err != nil {
		return err
	}
	*n = NoncurrentVersionTransition(val)
	n.set = true
	return nil
}

// IsNull returns true if NoncurrentTransition doesn't refer to any storage-class.
// Note: It supports immediate transition, i.e zero noncurrent days.
func (n NoncurrentVersionTransition) IsNull() bool {
	return n.StorageClass == ""
}

// Validate returns an error with wrong value
func (n NoncurrentVersionTransition) Validate() error {
	if !n.set {
		return nil
	}
	if n.StorageClass == "" {
		return errXMLNotWellFormed
	}
	return nil
}

// NextDue returns upcoming NoncurrentVersionTransition date for obj if
// applicable, returns false otherwise.
func (n NoncurrentVersionTransition) NextDue(obj ObjectOpts) (time.Time, bool) {
	if obj.IsLatest || n.StorageClass == "" {
		return time.Time{}, false
	}
	// Days == 0 indicates immediate tiering, i.e object is eligible for tiering since it became noncurrent.
	if n.NoncurrentDays == 0 {
		return obj.SuccessorModTime, true
	}
	return ExpectedExpiryTime(obj.SuccessorModTime, int(n.NoncurrentDays)), true
}
