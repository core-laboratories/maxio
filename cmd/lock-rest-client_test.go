// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"testing"

	"github.com/core-laboratories/maxio/internal/dsync"
)

// Tests lock rpc client.
func TestLockRESTlient(t *testing.T) {
	// These should not be connectable.
	endpoint, err := NewEndpoint("http://localhost:9876")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	endpointLocal, err := NewEndpoint("http://localhost:9012")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	endpointLocal.IsLocal = true

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	err = initGlobalLockGrid(ctx, []PoolEndpoints{{Endpoints: Endpoints{endpoint, endpointLocal}}})
	if err != nil {
		t.Fatal(err)
	}

	lkClient := newlockRESTClient(endpoint)
	if lkClient.IsOnline() {
		t.Fatalf("unexpected result. connection was online")
	}

	// Attempt all calls.
	_, err = lkClient.RLock(t.Context(), dsync.LockArgs{})
	if err == nil {
		t.Fatal("Expected for Rlock to fail")
	}

	_, err = lkClient.Lock(t.Context(), dsync.LockArgs{})
	if err == nil {
		t.Fatal("Expected for Lock to fail")
	}

	_, err = lkClient.RUnlock(t.Context(), dsync.LockArgs{})
	if err == nil {
		t.Fatal("Expected for RUnlock to fail")
	}

	_, err = lkClient.Unlock(t.Context(), dsync.LockArgs{})
	if err == nil {
		t.Fatal("Expected for Unlock to fail")
	}
}
