// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"reflect"
	"testing"
)

func TestServerConfigFile(t *testing.T) {
	for _, testcase := range []struct {
		config      string
		expectedErr bool
		hash        string
	}{
		{
			config:      "testdata/config/1.yaml",
			expectedErr: false,
			hash:        "hash:02bf70285dc71f76",
		},
		{
			config:      "testdata/config/2.yaml",
			expectedErr: false,
			hash:        "hash:676d2da00f71f205",
		},
		{
			config:      "testdata/config/invalid.yaml",
			expectedErr: true,
		},
		{
			config:      "testdata/config/invalid-types.yaml",
			expectedErr: true,
		},
		{
			config:      "testdata/config/invalid-disks.yaml",
			expectedErr: true,
		},
	} {
		t.Run(testcase.config, func(t *testing.T) {
			sctx := &serverCtxt{}
			err := mergeServerCtxtFromConfigFile(testcase.config, sctx)
			if testcase.expectedErr && err == nil {
				t.Error("expected failure, got success")
			}
			if !testcase.expectedErr && err != nil {
				t.Error("expected success, got failure", err)
			}
			if err == nil {
				if len(sctx.Layout.pools) != 2 {
					t.Error("expected parsed pools to be 2, not", len(sctx.Layout.pools))
				}
				if sctx.Layout.pools[0].cmdline != testcase.hash {
					t.Error("expected hash", testcase.hash, "got", sctx.Layout.pools[0].cmdline)
				}
			}
		})
	}
}

// Tests initializing new object layer.
func TestNewObjectLayer(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	// Tests for ErasureSD object layer.
	nDisks := 1
	disks, err := getRandomDisks(nDisks)
	if err != nil {
		t.Fatal("Failed to create drives for the backend")
	}
	defer removeRoots(disks)

	obj, err := newObjectLayer(ctx, mustGetPoolEndpoints(0, disks...))
	if err != nil {
		t.Fatal("Unexpected object layer initialization error", err)
	}

	_, ok := obj.(*erasureServerPools)
	if !ok {
		t.Fatal("Unexpected object layer detected", reflect.TypeOf(obj))
	}

	// Tests for Erasure object layer initialization.

	// Create temporary backend for the test server.
	nDisks = 16
	disks, err = getRandomDisks(nDisks)
	if err != nil {
		t.Fatal("Failed to create drives for the backend")
	}
	defer removeRoots(disks)

	obj, err = newObjectLayer(ctx, mustGetPoolEndpoints(0, disks...))
	if err != nil {
		t.Fatal("Unexpected object layer initialization error", err)
	}

	_, ok = obj.(*erasureServerPools)
	if !ok {
		t.Fatal("Unexpected object layer detected", reflect.TypeOf(obj))
	}
}
