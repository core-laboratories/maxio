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
	"reflect"
	"testing"
)

// Tests stripping standard ports from apiEndpoints.
func TestStripStandardPorts(t *testing.T) {
	apiEndpoints := []string{"http://127.0.0.1:9000", "http://127.0.0.2:80", "https://127.0.0.3:443"}
	expectedAPIEndpoints := []string{"http://127.0.0.1:9000", "http://127.0.0.2", "https://127.0.0.3"}
	newAPIEndpoints := stripStandardPorts(apiEndpoints, "")

	if !reflect.DeepEqual(expectedAPIEndpoints, newAPIEndpoints) {
		t.Fatalf("Expected %#v, got %#v", expectedAPIEndpoints, newAPIEndpoints)
	}

	apiEndpoints = []string{"http://%%%%%:9000"}
	newAPIEndpoints = stripStandardPorts(apiEndpoints, "")
	if !reflect.DeepEqual(apiEndpoints, newAPIEndpoints) {
		t.Fatalf("Expected %#v, got %#v", apiEndpoints, newAPIEndpoints)
	}

	apiEndpoints = []string{"http://127.0.0.1:443", "https://127.0.0.1:80"}
	newAPIEndpoints = stripStandardPorts(apiEndpoints, "")
	if !reflect.DeepEqual(apiEndpoints, newAPIEndpoints) {
		t.Fatalf("Expected %#v, got %#v", apiEndpoints, newAPIEndpoints)
	}
}

// Test printing server common message.
func TestPrintServerCommonMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	obj, fsDir, err := prepareFS(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fsDir)
	if err = newTestConfig(globalMinioDefaultRegion, obj); err != nil {
		t.Fatal(err)
	}

	apiEndpoints := []string{"http://127.0.0.1:9000"}
	printServerCommonMsg(apiEndpoints)
}

// Tests print cli access message.
func TestPrintCLIAccessMsg(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	obj, fsDir, err := prepareFS(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fsDir)
	if err = newTestConfig(globalMinioDefaultRegion, obj); err != nil {
		t.Fatal(err)
	}

	apiEndpoints := []string{"http://127.0.0.1:9000"}
	printCLIAccessMsg(apiEndpoints[0], "myminio")
}

// Test print startup message.
func TestPrintStartupMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	obj, fsDir, err := prepareFS(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fsDir)
	if err = newTestConfig(globalMinioDefaultRegion, obj); err != nil {
		t.Fatal(err)
	}

	apiEndpoints := []string{"http://127.0.0.1:9000"}
	printStartupMessage(apiEndpoints, nil)
}
