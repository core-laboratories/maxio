// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import (
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

	return "arn:minio:s3-object-lambda:" + arn.region + ":" + arn.TargetID.String()
}

// ParseARN - parses string to ARN.
func ParseARN(s string) (*ARN, error) {
	// ARN must be in the format of arn:minio:s3-object-lambda:<REGION>:<ID>:<TYPE>
	if !strings.HasPrefix(s, "arn:minio:s3-object-lambda:") {
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
