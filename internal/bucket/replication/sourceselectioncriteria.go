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

// ReplicaModifications specifies if replica modification sync is enabled
type ReplicaModifications struct {
	Status Status `xml:"Status" json:"Status"`
}

// SourceSelectionCriteria - specifies additional source selection criteria in ReplicationConfiguration.
type SourceSelectionCriteria struct {
	ReplicaModifications ReplicaModifications `xml:"ReplicaModifications" json:"ReplicaModifications"`
}

// IsValid - checks whether SourceSelectionCriteria is valid or not.
func (s SourceSelectionCriteria) IsValid() bool {
	return s.ReplicaModifications.Status == Enabled || s.ReplicaModifications.Status == Disabled
}

// Validate source selection criteria
func (s SourceSelectionCriteria) Validate() error {
	if (s == SourceSelectionCriteria{}) {
		return nil
	}
	if !s.IsValid() {
		return errInvalidSourceSelectionCriteria
	}
	return nil
}

// UnmarshalXML - decodes XML data.
func (s *SourceSelectionCriteria) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	// Make subtype to avoid recursive UnmarshalXML().
	type sourceSelectionCriteria SourceSelectionCriteria
	ssc := sourceSelectionCriteria{}
	if err := dec.DecodeElement(&ssc, &start); err != nil {
		return err
	}
	if len(ssc.ReplicaModifications.Status) == 0 {
		ssc.ReplicaModifications.Status = Enabled
	}
	*s = SourceSelectionCriteria(ssc)
	return nil
}

// MarshalXML - encodes to XML data.
func (s SourceSelectionCriteria) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if s.IsValid() {
		if err := e.EncodeElement(s.ReplicaModifications, xml.StartElement{Name: xml.Name{Local: "ReplicaModifications"}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
