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

func TestExtractPrefixAndSuffix(t *testing.T) {
	specs := []struct {
		path, prefix, suffix string
		expected             string
	}{
		{"config/iam/groups/foo.json", "config/iam/groups/", ".json", "foo"},
		{"config/iam/groups/./foo.json", "config/iam/groups/", ".json", "foo"},
		{"config/iam/groups/foo/config.json", "config/iam/groups/", "/config.json", "foo"},
		{"config/iam/groups/foo/config.json", "config/iam/groups/", "config.json", "foo"},
	}
	for i, test := range specs {
		result := extractPathPrefixAndSuffix(test.path, test.prefix, test.suffix)
		if result != test.expected {
			t.Errorf("unexpected result on test[%v]: expected[%s] but had [%s]", i, test.expected, result)
		}
	}
}
