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

var errInvalidDaysDelMarkerExpiration = Errorf("Days must be a positive integer with DelMarkerExpiration")

// DelMarkerExpiration used to xml encode/decode ILM action by the same name
type DelMarkerExpiration struct {
	XMLName xml.Name `xml:"DelMarkerExpiration"`
	Days    int      `xml:"Days,omitempty"`
}

// Empty returns if a DelMarkerExpiration XML element is empty.
// Used to detect if lifecycle.Rule contained a DelMarkerExpiration element.
func (de DelMarkerExpiration) Empty() bool {
	return de.Days == 0
}

// UnmarshalXML decodes a single XML element into a DelMarkerExpiration value
func (de *DelMarkerExpiration) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type delMarkerExpiration DelMarkerExpiration
	var dexp delMarkerExpiration
	err := dec.DecodeElement(&dexp, &start)
	if err != nil {
		return err
	}

	if dexp.Days <= 0 {
		return errInvalidDaysDelMarkerExpiration
	}

	*de = DelMarkerExpiration(dexp)
	return nil
}

// MarshalXML encodes a DelMarkerExpiration value into an XML element
func (de DelMarkerExpiration) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if de.Empty() {
		return nil
	}

	type delMarkerExpiration DelMarkerExpiration
	return enc.EncodeElement(delMarkerExpiration(de), start)
}

// NextDue returns upcoming DelMarkerExpiration date for obj if
// applicable, returns false otherwise.
func (de DelMarkerExpiration) NextDue(obj ObjectOpts) (time.Time, bool) {
	if !obj.IsLatest || !obj.DeleteMarker {
		return time.Time{}, false
	}

	return ExpectedExpiryTime(obj.ModTime, de.Days), true
}
