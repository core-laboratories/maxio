// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package arn

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ARN structure:
//
// arn:partition:service:region:account-id:resource-type/resource-id
//
// In this implementation, account-id is empty.
//
// Reference: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html

const (
	arnPrefixArn        = "arn"
	arnPartitionMinio   = "minio"
	arnServiceIAM       = "iam"
	arnResourceTypeRole = "role"
)

// ARN - representation of resources based on AWS ARNs.
type ARN struct {
	Partition    string
	Service      string
	Region       string
	ResourceType string
	ResourceID   string
}

// Allows english letters, numbers, '.', '-', '_' and '/'. Starts with a
// letter or digit. At least 1 character long.
var validResourceIDRegex = regexp.MustCompile(`[A-Za-z0-9_/\.-]+$`)

// NewIAMRoleARN - returns an ARN for a role in MinIO.
func NewIAMRoleARN(resourceID, serverRegion string) (ARN, error) {
	if !validResourceIDRegex.MatchString(resourceID) {
		return ARN{}, fmt.Errorf("invalid resource ID: %s", resourceID)
	}
	return ARN{
		Partition:    arnPartitionMinio,
		Service:      arnServiceIAM,
		Region:       serverRegion,
		ResourceType: arnResourceTypeRole,
		ResourceID:   resourceID,
	}, nil
}

// String - returns string representation of the ARN.
func (arn ARN) String() string {
	return strings.Join(
		[]string{
			arnPrefixArn,
			arn.Partition,
			arn.Service,
			arn.Region,
			"", // account-id is always empty in this implementation
			arn.ResourceType + "/" + arn.ResourceID,
		},
		":",
	)
}

// Parse - parses an ARN string into a type.
func Parse(arnStr string) (arn ARN, err error) {
	ps := strings.Split(arnStr, ":")
	if len(ps) != 6 || ps[0] != string(arnPrefixArn) {
		err = errors.New("invalid ARN string format")
		return arn, err
	}

	if ps[1] != string(arnPartitionMinio) {
		err = errors.New("invalid ARN - bad partition field")
		return arn, err
	}

	if ps[2] != string(arnServiceIAM) {
		err = errors.New("invalid ARN - bad service field")
		return arn, err
	}

	// ps[3] is region and is not validated here. If the region is invalid,
	// the ARN would not match any configured ARNs in the server.
	if ps[4] != "" {
		err = errors.New("invalid ARN - unsupported account-id field")
		return arn, err
	}

	res := strings.SplitN(ps[5], "/", 2)
	if len(res) != 2 {
		err = errors.New("invalid ARN - resource does not contain a \"/\"")
		return arn, err
	}

	if res[0] != string(arnResourceTypeRole) {
		err = errors.New("invalid ARN: resource type is invalid")
		return arn, err
	}

	if !validResourceIDRegex.MatchString(res[1]) {
		err = fmt.Errorf("invalid resource ID: %s", res[1])
		return arn, err
	}

	arn = ARN{
		Partition:    arnPartitionMinio,
		Service:      arnServiceIAM,
		Region:       ps[3],
		ResourceType: arnResourceTypeRole,
		ResourceID:   res[1],
	}
	return arn, err
}
