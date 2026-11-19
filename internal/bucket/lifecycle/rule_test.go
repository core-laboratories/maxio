// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lifecycle

import (
	"encoding/xml"
	"fmt"
	"testing"
)

// TestInvalidRules checks if Rule xml with invalid elements returns
// appropriate errors on validation
func TestInvalidRules(t *testing.T) {
	invalidTestCases := []struct {
		inputXML    string
		expectedErr error
	}{
		{ // Rule with ID longer than 255 characters
			inputXML: ` <Rule>
	                    <ID> babababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababab </ID>
	                    </Rule>`,
			expectedErr: errInvalidRuleID,
		},
		{ // Rule with empty ID
			inputXML: `<Rule>
							<ID></ID>
							<Filter><Prefix></Prefix></Filter>
							<Expiration>
								<Days>365</Days>
							</Expiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: nil,
		},
		{ // Rule with empty status
			inputXML: ` <Rule>
			                  <ID>rule with empty status</ID>
                              <Status></Status>
	                    </Rule>`,
			expectedErr: errEmptyRuleStatus,
		},
		{ // Rule with invalid status
			inputXML: ` <Rule>
			                  <ID>rule with invalid status</ID>
                              <Status>OK</Status>
	                    </Rule>`,
			expectedErr: errInvalidRuleStatus,
		},
		{ // Rule with negative values for ObjectSizeLessThan
			inputXML: `<Rule>
				<ID>negative-obj-size-less-than</ID>
				<Filter><ObjectSizeLessThan>-1</ObjectSizeLessThan></Filter>
				<Expiration>
					<Days>365</Days>
				</Expiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: errXMLNotWellFormed,
		},
		{ // Rule with negative values for And>ObjectSizeLessThan
			inputXML: `<Rule>
				<ID>negative-and-obj-size-less-than</ID>
				<Filter><And><ObjectSizeLessThan>-1</ObjectSizeLessThan></And></Filter>
				<Expiration>
					<Days>365</Days>
				</Expiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: errXMLNotWellFormed,
		},
		{ // Rule with negative values for ObjectSizeGreaterThan
			inputXML: `<Rule>
				<ID>negative-obj-size-greater-than</ID>
				<Filter><ObjectSizeGreaterThan>-1</ObjectSizeGreaterThan></Filter>
				<Expiration>
					<Days>365</Days>
				</Expiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: errXMLNotWellFormed,
		},
		{ // Rule with negative values for And>ObjectSizeGreaterThan
			inputXML: `<Rule>
				<ID>negative-and-obj-size-greater-than</ID>
				<Filter><And><ObjectSizeGreaterThan>-1</ObjectSizeGreaterThan></And></Filter>
				<Expiration>
					<Days>365</Days>
				</Expiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: errXMLNotWellFormed,
		},
		{
			inputXML: `<Rule>
				<ID>Rule with a tag and DelMarkerExpiration</ID>
				<Filter><Tag><Key>k1</Key><Value>v1</Value></Tag></Filter>
				<DelMarkerExpiration>
					<Days>365</Days>
				</DelMarkerExpiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: errInvalidRuleDelMarkerExpiration,
		},
		{
			inputXML: `<Rule>
				<ID>Rule with multiple tags and DelMarkerExpiration</ID>
				<Filter><And>
				<Tag><Key>k1</Key><Value>v1</Value></Tag>
				<Tag><Key>k2</Key><Value>v2</Value></Tag>
				</And></Filter>
				<DelMarkerExpiration>
					<Days>365</Days>
				</DelMarkerExpiration>
                            <Status>Enabled</Status>
	                    </Rule>`,
			expectedErr: errInvalidRuleDelMarkerExpiration,
		},
	}

	for i, tc := range invalidTestCases {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			var rule Rule
			err := xml.Unmarshal([]byte(tc.inputXML), &rule)
			if err != nil {
				t.Fatal(err)
			}

			if err := rule.Validate(); err != tc.expectedErr {
				t.Fatalf("%d: Expected %v but got %v", i+1, tc.expectedErr, err)
			}
		})
	}
}
