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
	"sync"
	"testing"

	"github.com/core-laboratories/maxio/internal/dsync"
)

// Helper function to create a lock server for testing
func createLockTestServer(ctx context.Context, t *testing.T) (string, *lockRESTServer, string) {
	obj, fsDir, err := prepareFS(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err = newTestConfig(globalMinioDefaultRegion, obj); err != nil {
		t.Fatalf("unable initialize config file, %s", err)
	}

	locker := &lockRESTServer{
		ll: &localLocker{
			mutex:   sync.Mutex{},
			lockMap: make(map[string][]lockRequesterInfo),
		},
	}
	creds := globalActiveCred
	token, err := authenticateNode(creds.AccessKey, creds.SecretKey)
	if err != nil {
		t.Fatal(err)
	}
	return fsDir, locker, token
}

// Test function to remove lock entries from map based on name & uid combination
func TestLockRpcServerRemoveEntry(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	testPath, locker, _ := createLockTestServer(ctx, t)
	defer os.RemoveAll(testPath)

	lockRequesterInfo1 := lockRequesterInfo{
		Owner:           "owner",
		Writer:          true,
		UID:             "0123-4567",
		Timestamp:       UTCNow().UnixNano(),
		TimeLastRefresh: UTCNow().UnixNano(),
	}
	lockRequesterInfo2 := lockRequesterInfo{
		Owner:           "owner",
		Writer:          true,
		UID:             "89ab-cdef",
		Timestamp:       UTCNow().UnixNano(),
		TimeLastRefresh: UTCNow().UnixNano(),
	}

	locker.ll.lockMap["name"] = []lockRequesterInfo{
		lockRequesterInfo1,
		lockRequesterInfo2,
	}

	lri := locker.ll.lockMap["name"]

	// test unknown uid
	if locker.ll.removeEntry("name", dsync.LockArgs{
		Owner: "owner",
		UID:   "unknown-uid",
	}, &lri) {
		t.Errorf("Expected %#v, got %#v", false, true)
	}

	if !locker.ll.removeEntry("name", dsync.LockArgs{
		Owner: "owner",
		UID:   "0123-4567",
	}, &lri) {
		t.Errorf("Expected %#v, got %#v", true, false)
	} else {
		gotLri := locker.ll.lockMap["name"]
		expectedLri := []lockRequesterInfo{lockRequesterInfo2}
		if !reflect.DeepEqual(expectedLri, gotLri) {
			t.Errorf("Expected %#v, got %#v", expectedLri, gotLri)
		}
	}

	if !locker.ll.removeEntry("name", dsync.LockArgs{
		Owner: "owner",
		UID:   "89ab-cdef",
	}, &lri) {
		t.Errorf("Expected %#v, got %#v", true, false)
	} else {
		gotLri := locker.ll.lockMap["name"]
		expectedLri := []lockRequesterInfo(nil)
		if !reflect.DeepEqual(expectedLri, gotLri) {
			t.Errorf("Expected %#v, got %#v", expectedLri, gotLri)
		}
	}
}
