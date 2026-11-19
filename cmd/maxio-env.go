// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"os"
	"strings"
)

// applyMaxIOEnvAliases lets the MaxIO configuration names override their
// inherited MINIO counterparts. Native MaxIO IPFS settings need no alias.
func applyMaxIOEnvAliases() {
	for _, entry := range os.Environ() {
		name, value, ok := strings.Cut(entry, "=")
		if !ok || !strings.HasPrefix(name, "MAXIO_") ||
			strings.HasPrefix(name, "MAXIO_IPFS_") || name == "MAXIO_RELEASE_TAG" {
			continue
		}
		legacyName := "MINIO_" + strings.TrimPrefix(name, "MAXIO_")
		_ = os.Setenv(legacyName, value)
	}
}
