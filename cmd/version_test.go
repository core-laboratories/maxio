// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	Version = "2017-05-07T06:37:49Z"
	_, err := time.Parse(time.RFC3339, Version)
	if err != nil {
		t.Fatal(err)
	}
}
