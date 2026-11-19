// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"testing"
)

func TestSplitPath(t *testing.T) {
	cases := []struct {
		path                          string
		secondIndex                   bool
		expectedListKey, expectedItem string
	}{
		{"format.json", false, "format.json", ""},
		{"users/tester.json", false, "users/", "tester.json"},
		{"groups/test/group.json", false, "groups/", "test/group.json"},
		{"policydb/groups/testgroup.json", true, "policydb/groups/", "testgroup.json"},
		{
			"policydb/sts-users/uid=slash/user,ou=people,ou=swengg,dc=min,dc=io.json", true,
			"policydb/sts-users/", "uid=slash/user,ou=people,ou=swengg,dc=min,dc=io.json",
		},
		{
			"policydb/sts-users/uid=slash/user/twice,ou=people,ou=swengg,dc=min,dc=io.json", true,
			"policydb/sts-users/", "uid=slash/user/twice,ou=people,ou=swengg,dc=min,dc=io.json",
		},
		{
			"policydb/groups/cn=project/d,ou=groups,ou=swengg,dc=min,dc=io.json", true,
			"policydb/groups/", "cn=project/d,ou=groups,ou=swengg,dc=min,dc=io.json",
		},
	}
	for i, test := range cases {
		listKey, item := splitPath(test.path, test.secondIndex)
		if listKey != test.expectedListKey || item != test.expectedItem {
			t.Errorf("unexpected result on test[%v]: expected[%s, %s] but got [%s, %s]", i, test.expectedListKey, test.expectedItem, listKey, item)
		}
	}
}
