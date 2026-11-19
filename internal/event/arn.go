// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import (
	"encoding/xml"
	"strings"
)

// ARN - SQS resource name representation.
type ARN struct {
	TargetID
	region string
}

// String - returns string representation.
func (arn ARN) String() string {
	if arn.ID == "" && arn.Name == "" && arn.region == "" {
		return ""
	}

	return "arn:minio:sqs:" + arn.region + ":" + arn.TargetID.String()
}

// MarshalXML - encodes to XML data.
func (arn ARN) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(arn.String(), start)
}

// UnmarshalXML - decodes XML data.
func (arn *ARN) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	parsedARN, err := parseARN(s)
	if err != nil {
		return err
	}

	*arn = *parsedARN
	return nil
}

// parseARN - parses string to ARN.
func parseARN(s string) (*ARN, error) {
	// ARN must be in the format of arn:minio:sqs:<REGION>:<ID>:<TYPE>
	if !strings.HasPrefix(s, "arn:minio:sqs:") {
		return nil, &ErrInvalidARN{s}
	}

	tokens := strings.Split(s, ":")
	if len(tokens) != 6 {
		return nil, &ErrInvalidARN{s}
	}

	if tokens[4] == "" || tokens[5] == "" {
		return nil, &ErrInvalidARN{s}
	}

	return &ARN{
		region: tokens[3],
		TargetID: TargetID{
			ID:   tokens[4],
			Name: tokens[5],
		},
	}, nil
}
