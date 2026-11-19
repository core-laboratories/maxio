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

// Tests - mkdirAll()
func TestOSMkdirAll(t *testing.T) {
	// create xlStorage test setup
	_, path, err := newXLStorageTestSetup(t)
	if err != nil {
		t.Fatalf("Unable to create xlStorage test setup, %s", err)
	}

	if err = mkdirAll("", 0o777, ""); err != errInvalidArgument {
		t.Fatal("Unexpected error", err)
	}

	if err = mkdirAll(pathJoin(path, "my-obj-del-0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"), 0o777, ""); err != errFileNameTooLong {
		t.Fatal("Unexpected error", err)
	}

	if err = mkdirAll(pathJoin(path, "success-vol", "success-object"), 0o777, ""); err != nil {
		t.Fatal("Unexpected error", err)
	}
}

// Tests - renameAll()
func TestOSRenameAll(t *testing.T) {
	// create xlStorage test setup
	_, path, err := newXLStorageTestSetup(t)
	if err != nil {
		t.Fatalf("Unable to create xlStorage test setup, %s", err)
	}

	if err = mkdirAll(pathJoin(path, "testvolume1"), 0o777, ""); err != nil {
		t.Fatal(err)
	}
	if err = renameAll("", "foo", ""); err != errInvalidArgument {
		t.Fatal(err)
	}
	if err = renameAll("foo", "", ""); err != errInvalidArgument {
		t.Fatal(err)
	}
	if err = renameAll(pathJoin(path, "testvolume1"), pathJoin(path, "testvolume2"), ""); err != nil {
		t.Fatal(err)
	}
	if err = renameAll(pathJoin(path, "testvolume1"), pathJoin(path, "testvolume2"), ""); err != errFileNotFound {
		t.Fatal(err)
	}
	if err = renameAll(pathJoin(path, "my-obj-del-0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"), pathJoin(path, "testvolume2"), ""); err != errFileNameTooLong {
		t.Fatal("Unexpected error", err)
	}
	if err = renameAll(pathJoin(path, "testvolume1"), pathJoin(path, "my-obj-del-0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"), ""); err != errFileNameTooLong {
		t.Fatal("Unexpected error", err)
	}
}
