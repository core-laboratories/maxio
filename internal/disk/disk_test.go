//go:build !netbsd && !solaris
// +build !netbsd,!solaris

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package disk_test

import (
	"testing"

	"github.com/core-laboratories/maxio/internal/disk"
)

func TestFree(t *testing.T) {
	di, err := disk.GetInfo(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}

	if di.FSType == "UNKNOWN" {
		t.Error("Unexpected FSType", di.FSType)
	}
}
