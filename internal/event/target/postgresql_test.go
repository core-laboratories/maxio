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

// TestPostgreSQLRegistration checks if postgres driver
// is registered and fails otherwise.
func TestPostgreSQLRegistration(t *testing.T) {
	var found bool
	if slices.Contains(sql.Drivers(), "postgres") {
		found = true
	}
	if !found {
		t.Fatal("postgres driver not registered")
	}
}

func TestPsqlTableNameValidation(t *testing.T) {
	validTables := []string{"táblë", "table", "TableName", "\"Table name\"", "\"✅✅\"", "table$one", "\"táblë\""}
	invalidTables := []string{"table name", "table \"name\"", "✅✅", "$table$"}

	for _, name := range validTables {
		if err := validatePsqlTableName(name); err != nil {
			t.Errorf("Should be valid: %s - %s", name, err)
		}
	}
	for _, name := range invalidTables {
		if err := validatePsqlTableName(name); err != errInvalidPsqlTablename {
			t.Errorf("Should be invalid: %s - %s", name, err)
		}
	}
}
