// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lifecycle

import (
	"encoding/xml"
	"io"
	"unicode/utf8"
)

// Tag - a tag for a lifecycle configuration Rule filter.
type Tag struct {
	XMLName xml.Name `xml:"Tag"`
	Key     string   `xml:"Key,omitempty"`
	Value   string   `xml:"Value,omitempty"`
}

var (
	errInvalidTagKey   = Errorf("The TagKey you have provided is invalid")
	errInvalidTagValue = Errorf("The TagValue you have provided is invalid")

	errDuplicatedXMLTag = Errorf("duplicated XML Tag")
	errUnknownXMLTag    = Errorf("unknown XML Tag")
)

// UnmarshalXML - decodes XML data.
func (tag *Tag) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var keyAlreadyParsed, valueAlreadyParsed bool
	for {
		// Read tokens from the XML document in a stream.
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if se, ok := t.(xml.StartElement); ok {
			var s string
			if err = d.DecodeElement(&s, &se); err != nil {
				return err
			}
			switch se.Name.Local {
			case "Key":
				if keyAlreadyParsed {
					return errDuplicatedXMLTag
				}
				tag.Key = s
				keyAlreadyParsed = true
			case "Value":
				if valueAlreadyParsed {
					return errDuplicatedXMLTag
				}
				tag.Value = s
				valueAlreadyParsed = true
			default:
				return errUnknownXMLTag
			}
		}
	}

	return nil
}

func (tag Tag) String() string {
	return tag.Key + "=" + tag.Value
}

// IsEmpty returns whether this tag is empty or not.
func (tag Tag) IsEmpty() bool {
	return tag.Key == ""
}

// Validate checks this tag.
func (tag Tag) Validate() error {
	if len(tag.Key) == 0 || utf8.RuneCountInString(tag.Key) > 128 {
		return errInvalidTagKey
	}

	if utf8.RuneCountInString(tag.Value) > 256 {
		return errInvalidTagValue
	}

	return nil
}
