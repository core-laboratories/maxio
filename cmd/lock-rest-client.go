// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"errors"

	"github.com/core-laboratories/maxio/internal/dsync"
	"github.com/core-laboratories/maxio/internal/grid"
)

// lockRESTClient is authenticable lock REST client
type lockRESTClient struct {
	connection *grid.Connection
}

// IsOnline - returns whether REST client failed to connect or not.
func (c *lockRESTClient) IsOnline() bool {
	return c.connection.State() == grid.StateConnected
}

// Not a local locker
func (c *lockRESTClient) IsLocal() bool {
	return false
}

// Close - marks the client as closed.
func (c *lockRESTClient) Close() error {
	return nil
}

// String - returns the remote host of the connection.
func (c *lockRESTClient) String() string {
	return c.connection.Remote
}

func (c *lockRESTClient) call(ctx context.Context, h *grid.SingleHandler[*dsync.LockArgs, *dsync.LockResp], args *dsync.LockArgs) (ok bool, err error) {
	r, err := h.Call(ctx, c.connection, args)
	if err != nil {
		return false, err
	}
	defer h.PutResponse(r)
	ok = r.Code == dsync.RespOK
	switch r.Code {
	case dsync.RespLockConflict, dsync.RespLockNotFound, dsync.RespOK:
	// no error
	case dsync.RespLockNotInitialized:
		err = errLockNotInitialized
	default:
		err = errors.New(r.Err)
	}
	return ok, err
}

// RLock calls read lock REST API.
func (c *lockRESTClient) RLock(ctx context.Context, args dsync.LockArgs) (reply bool, err error) {
	return c.call(ctx, lockRPCRLock, &args)
}

// Lock calls lock REST API.
func (c *lockRESTClient) Lock(ctx context.Context, args dsync.LockArgs) (reply bool, err error) {
	return c.call(ctx, lockRPCLock, &args)
}

// RUnlock calls read unlock REST API.
func (c *lockRESTClient) RUnlock(ctx context.Context, args dsync.LockArgs) (reply bool, err error) {
	return c.call(ctx, lockRPCRUnlock, &args)
}

// Refresh calls Refresh REST API.
func (c *lockRESTClient) Refresh(ctx context.Context, args dsync.LockArgs) (reply bool, err error) {
	return c.call(ctx, lockRPCRefresh, &args)
}

// Unlock calls write unlock RPC.
func (c *lockRESTClient) Unlock(ctx context.Context, args dsync.LockArgs) (reply bool, err error) {
	return c.call(ctx, lockRPCUnlock, &args)
}

// ForceUnlock calls force unlock handler to forcibly unlock an active lock.
func (c *lockRESTClient) ForceUnlock(ctx context.Context, args dsync.LockArgs) (reply bool, err error) {
	return c.call(ctx, lockRPCForceUnlock, &args)
}

func newLockAPI(endpoint Endpoint) dsync.NetLocker {
	if endpoint.IsLocal {
		return globalLockServer
	}
	return newlockRESTClient(endpoint)
}

// Returns a lock rest client.
func newlockRESTClient(ep Endpoint) *lockRESTClient {
	return &lockRESTClient{globalLockGrid.Load().Connection(ep.GridHost())}
}
