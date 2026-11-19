// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TargetID - holds identification and name strings of notification target.
type TargetID struct {
	ID   string
	Name string
}

// String - returns string representation.
func (tid TargetID) String() string {
	return tid.ID + ":" + tid.Name
}

// ToARN - converts to ARN.
func (tid TargetID) ToARN(region string) ARN {
	return ARN{TargetID: tid, region: region}
}

// MarshalJSON - encodes to JSON data.
func (tid TargetID) MarshalJSON() ([]byte, error) {
	return json.Marshal(tid.String())
}

// UnmarshalJSON - decodes JSON data.
func (tid *TargetID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	targetID, err := parseTargetID(s)
	if err != nil {
		return err
	}

	*tid = *targetID
	return nil
}

// parseTargetID - parses string to TargetID.
func parseTargetID(s string) (*TargetID, error) {
	tokens := strings.Split(s, ":")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("invalid TargetID format '%v'", s)
	}

	return &TargetID{
		ID:   tokens[0],
		Name: tokens[1],
	}, nil
}
