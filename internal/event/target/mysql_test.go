// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package target

import (
	"database/sql"
	"slices"
	"testing"
)

// TestPostgreSQLRegistration checks if sql driver
// is registered and fails otherwise.
func TestMySQLRegistration(t *testing.T) {
	var found bool
	if slices.Contains(sql.Drivers(), "mysql") {
		found = true
	}
	if !found {
		t.Fatal("mysql driver not registered")
	}
}
