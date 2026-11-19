// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package csv

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	none   = "none"
	use    = "use"
	ignore = "ignore"

	defaultRecordDelimiter      = "\n"
	defaultFieldDelimiter       = ","
	defaultQuoteCharacter       = `"`
	defaultQuoteEscapeCharacter = `"`
	defaultCommentCharacter     = "#"

	asneeded = "asneeded"
)

// ReaderArgs - represents elements inside <InputSerialization><CSV> in request XML.
type ReaderArgs struct {
	FileHeaderInfo             string `xml:"FileHeaderInfo"`
	RecordDelimiter            string `xml:"RecordDelimiter"`
	FieldDelimiter             string `xml:"FieldDelimiter"`
	QuoteCharacter             string `xml:"QuoteCharacter"`
	QuoteEscapeCharacter       string `xml:"QuoteEscapeCharacter"`
	CommentCharacter           string `xml:"Comments"`
	AllowQuotedRecordDelimiter bool   `xml:"AllowQuotedRecordDelimiter"`
	unmarshaled                bool
}

// IsEmpty - returns whether reader args is empty or not.
func (args *ReaderArgs) IsEmpty() bool {
	return !args.unmarshaled
}

// UnmarshalXML - decodes XML data.
func (args *ReaderArgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	args.FileHeaderInfo = none
	args.RecordDelimiter = defaultRecordDelimiter
	args.FieldDelimiter = defaultFieldDelimiter
	args.QuoteCharacter = defaultQuoteCharacter
	args.QuoteEscapeCharacter = defaultQuoteEscapeCharacter
	args.CommentCharacter = defaultCommentCharacter
	args.AllowQuotedRecordDelimiter = false

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
			tagName := se.Name.Local
			switch tagName {
			case "AllowQuotedRecordDelimiter":
				var b bool
				if err = d.DecodeElement(&b, &se); err != nil {
					return err
				}
				args.AllowQuotedRecordDelimiter = b
			default:
				var s string
				if err = d.DecodeElement(&s, &se); err != nil {
					return err
				}
				switch tagName {
				case "FileHeaderInfo":
					s = strings.ToLower(s)
					if len(s) != 0 {
						if s != none && s != use && s != ignore {
							return errors.New("unsupported FileHeaderInfo")
						}
						args.FileHeaderInfo = s
					}

				case "RecordDelimiter":
					if len(s) != 0 {
						args.RecordDelimiter = s
					}
				case "FieldDelimiter":
					if len(s) != 0 {
						args.FieldDelimiter = s
					}
				case "QuoteCharacter":
					if utf8.RuneCountInString(s) > 1 {
						return fmt.Errorf("unsupported QuoteCharacter '%v'", s)
					}
					args.QuoteCharacter = s
				case "QuoteEscapeCharacter":
					switch utf8.RuneCountInString(s) {
					case 0:
						args.QuoteEscapeCharacter = defaultQuoteEscapeCharacter
					case 1:
						args.QuoteEscapeCharacter = s
					default:
						return fmt.Errorf("unsupported QuoteEscapeCharacter '%v'", s)
					}
				case "Comments":
					if len(s) != 0 {
						args.CommentCharacter = s
					}
				default:
					return errors.New("unrecognized option")
				}
			}
		}
	}

	args.unmarshaled = true
	return nil
}

// WriterArgs - represents elements inside <OutputSerialization><CSV/> in request XML.
type WriterArgs struct {
	QuoteFields          string `xml:"QuoteFields"`
	RecordDelimiter      string `xml:"RecordDelimiter"`
	FieldDelimiter       string `xml:"FieldDelimiter"`
	QuoteCharacter       string `xml:"QuoteCharacter"`
	QuoteEscapeCharacter string `xml:"QuoteEscapeCharacter"`
	unmarshaled          bool
}

// IsEmpty - returns whether writer args is empty or not.
func (args *WriterArgs) IsEmpty() bool {
	return !args.unmarshaled
}

// UnmarshalXML - decodes XML data.
func (args *WriterArgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	args.QuoteFields = asneeded
	args.RecordDelimiter = defaultRecordDelimiter
	args.FieldDelimiter = defaultFieldDelimiter
	args.QuoteCharacter = defaultQuoteCharacter
	args.QuoteEscapeCharacter = defaultQuoteCharacter

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
			case "QuoteFields":
				args.QuoteFields = strings.ToLower(s)
			case "RecordDelimiter":
				args.RecordDelimiter = s
			case "FieldDelimiter":
				args.FieldDelimiter = s
			case "QuoteCharacter":
				switch utf8.RuneCountInString(s) {
				case 0:
					args.QuoteCharacter = "\x00"
				case 1:
					args.QuoteCharacter = s
				default:
					return fmt.Errorf("unsupported QuoteCharacter '%v'", s)
				}
			case "QuoteEscapeCharacter":
				switch utf8.RuneCountInString(s) {
				case 0:
					args.QuoteEscapeCharacter = defaultQuoteEscapeCharacter
				case 1:
					args.QuoteEscapeCharacter = s
				default:
					return fmt.Errorf("unsupported QuoteCharacter '%v'", s)
				}
			default:
				return errors.New("unrecognized option")
			}
		}
	}

	args.unmarshaled = true
	return nil
}
