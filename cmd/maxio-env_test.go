// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"os"
	"testing"
)

func TestApplyMaxIOEnvAliases(t *testing.T) {
	t.Setenv("MINIO_ROOT_USER", "legacy")
	t.Setenv("MAXIO_ROOT_USER", "maxio")
	t.Setenv("MINIO_IPFS_ENABLED", "legacy")
	t.Setenv("MAXIO_IPFS_ENABLED", "false")

	applyMaxIOEnvAliases()

	if got := os.Getenv("MINIO_ROOT_USER"); got != "maxio" {
		t.Fatalf("root user alias = %q, want maxio", got)
	}
	if got := os.Getenv("MINIO_IPFS_ENABLED"); got != "legacy" {
		t.Fatalf("native IPFS setting should not create an alias: got %q", got)
	}
}
