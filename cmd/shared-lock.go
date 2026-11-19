// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"time"
)

var sharedLockTimeout = newDynamicTimeoutWithOpts(dynamicTimeoutOpts{
	timeout:       30 * time.Second,
	minimum:       10 * time.Second,
	retryInterval: time.Minute,
})

type sharedLock struct {
	lockContext chan LockContext
}

func (ld sharedLock) backgroundRoutine(ctx context.Context, objAPI ObjectLayer, lockName string) {
	for {
		locker := objAPI.NewNSLock(minioMetaBucket, lockName)
		lkctx, err := locker.GetLock(ctx, sharedLockTimeout)
		if err != nil {
			continue
		}

	keepLock:
		for {
			select {
			case <-ctx.Done():
				return
			case <-lkctx.Context().Done():
				// The context of the lock is canceled, this can happen
				// if one lock lost quorum due to cluster instability
				// in that case, try to lock again.
				break keepLock
			case ld.lockContext <- lkctx:
				// Send the lock context to anyone asking for it
			}
		}
	}
}

func mergeContext(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx1.Done():
		case <-ctx2.Done():
		// The lock acquirer decides to cancel, exit this goroutine
		case <-ctx.Done():
		}

		cancel()
	}()
	return ctx, cancel
}

func (ld sharedLock) GetLock(ctx context.Context) (context.Context, context.CancelFunc) {
	l := <-ld.lockContext
	return mergeContext(l.Context(), ctx)
}

func newSharedLock(ctx context.Context, objAPI ObjectLayer, lockName string) *sharedLock {
	l := &sharedLock{
		lockContext: make(chan LockContext),
	}
	go l.backgroundRoutine(ctx, objAPI, lockName)

	return l
}
