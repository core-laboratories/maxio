// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"testing"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7/pkg/set"
)

// TestGetMissingSiteNames
func TestGetMissingSiteNames(t *testing.T) {
	testCases := []struct {
		currSites []madmin.PeerInfo
		oldDepIDs set.StringSet
		newDepIDs set.StringSet
		expNames  []string
	}{
		// Test1: missing some sites in replicated setup
		{
			[]madmin.PeerInfo{
				{Endpoint: "minio1:9000", Name: "minio1", DeploymentID: "dep1"},
				{Endpoint: "minio2:9000", Name: "minio2", DeploymentID: "dep2"},
				{Endpoint: "minio3:9000", Name: "minio3", DeploymentID: "dep3"},
			},
			set.CreateStringSet("dep1", "dep2", "dep3"),
			set.CreateStringSet("dep1"),
			[]string{"minio2", "minio3"},
		},
		// Test2: new site added that is not in replicated setup
		{
			[]madmin.PeerInfo{{Endpoint: "minio1:9000", Name: "minio1", DeploymentID: "dep1"}, {Endpoint: "minio2:9000", Name: "minio2", DeploymentID: "dep2"}, {Endpoint: "minio3:9000", Name: "minio3", DeploymentID: "dep3"}},
			set.CreateStringSet("dep1", "dep2", "dep3"),
			set.CreateStringSet("dep1", "dep2", "dep3", "dep4"),
			[]string{},
		},
		// Test3: not currently under site replication.
		{
			[]madmin.PeerInfo{},
			set.CreateStringSet(),
			set.CreateStringSet("dep1", "dep2", "dep3", "dep4"),
			[]string{},
		},
	}

	for i, tc := range testCases {
		names := getMissingSiteNames(tc.oldDepIDs, tc.newDepIDs, tc.currSites)
		if len(names) != len(tc.expNames) {
			t.Errorf("Test %d: Expected `%v`, got `%v`", i+1, tc.expNames, names)
		}
	}
}
