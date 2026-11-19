//go:build windows
// +build windows

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

// Test if various paths work as expected when converted to UNC form
func TestUNCPaths(t *testing.T) {
	testCases := []struct {
		objName string
		pass    bool
	}{
		{"/abcdef", true},
		{"/a/b/c/d/e/f/g", true},
		{string(bytes.Repeat([]byte("界"), 85)), true},
		// Each path component must be <= 255 bytes long.
		{string(bytes.Repeat([]byte("界"), 280)), false},
		{`/p/q/r/s/t`, true},
	}
	dir := t.TempDir()

	// Instantiate posix object to manage a disk
	fs, err := newLocalXLStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Create volume to use in conjunction with other StorageAPI's file API(s)
	err = fs.MakeVol(context.Background(), "voldir")
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err = fs.AppendFile(context.Background(), "voldir", test.objName, []byte("hello"))
			if err != nil && test.pass {
				t.Error(err)
			} else if err == nil && !test.pass {
				t.Error(err)
			}
			fs.Delete(context.Background(), "voldir", test.objName, DeleteOptions{
				Recursive: false,
				Immediate: false,
			})
		})
	}
}

// Test to validate xlStorage behavior on windows when a non-final path component is a file.
func TestUNCPathENOTDIR(t *testing.T) {
	// Instantiate posix object to manage a disk
	dir := t.TempDir()

	fs, err := newLocalXLStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Create volume to use in conjunction with other StorageAPI's file API(s)
	err = fs.MakeVol(context.Background(), "voldir")
	if err != nil {
		t.Fatal(err)
	}

	err = fs.AppendFile(context.Background(), "voldir", "/file", []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	// Try to create a file that includes a file in its path components.
	// In *nix, this returns syscall.ENOTDIR while in windows we receive the following error.
	err = fs.AppendFile(context.Background(), "voldir", "/file/obj1", []byte("hello"))
	if err != errFileAccessDenied {
		t.Errorf("expected: %s, got: %s", errFileAccessDenied, err)
	}
}
