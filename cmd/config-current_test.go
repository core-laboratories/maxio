// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"os"
	"testing"

	"github.com/core-laboratories/maxio/internal/config"
)

func TestServerConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	objLayer, fsDir, err := prepareFS(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fsDir)

	if err = newTestConfig(globalMinioDefaultRegion, objLayer); err != nil {
		t.Fatalf("Init Test config failed")
	}

	if globalSite.Region() != globalMinioDefaultRegion {
		t.Errorf("Expecting region `us-east-1` found %s", globalSite.Region())
	}

	// Set new region and verify.
	config.SetRegion(globalServerConfig, "us-west-1")
	site, err := config.LookupSite(
		globalServerConfig[config.SiteSubSys][config.Default],
		globalServerConfig[config.RegionSubSys][config.Default],
	)
	if err != nil {
		t.Fatal(err)
	}
	if site.Region() != "us-west-1" {
		t.Errorf("Expecting region `us-west-1` found %s", globalSite.Region())
	}

	if err := saveServerConfig(t.Context(), objLayer, globalServerConfig); err != nil {
		t.Fatalf("Unable to save updated config file %s", err)
	}

	// Initialize server config.
	if err := loadConfig(objLayer, nil); err != nil {
		t.Fatalf("Unable to initialize from updated config file %s", err)
	}
}
