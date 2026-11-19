// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// Tests getSource().
func TestGetSource(t *testing.T) {
	currentSource := func() string { return getSource(2) }
	_, file, line, _ := runtime.Caller(0)
	gotSource := currentSource()
	expectedSource := fmt.Sprintf("[%s:%d:TestGetSource()]", filepath.Base(file), line+1)
	if gotSource != expectedSource {
		t.Errorf("expected : %s, got : %s", expectedSource, gotSource)
	}
}

// Test lock race
func TestNSLockRace(t *testing.T) {
	t.Skip("long test skip it")

	ctx := t.Context()

	for i := range 10000 {
		nsLk := newNSLock(false)

		// lk1; ref=1
		if !nsLk.lock(ctx, "volume", "path", "source", "opsID", false, time.Second) {
			t.Fatal("failed to acquire lock")
		}

		// lk2
		lk2ch := make(chan struct{})
		go func() {
			defer close(lk2ch)
			nsLk.lock(ctx, "volume", "path", "source", "opsID", false, 1*time.Millisecond)
		}()
		time.Sleep(1 * time.Millisecond) // wait for goroutine to advance; ref=2

		// Unlock the 1st lock; ref=1 after this line
		nsLk.unlock("volume", "path", false)

		// Taking another lockMapMutex here allows queuing up additional lockers. This should
		// not be required but makes reproduction much easier.
		nsLk.lockMapMutex.Lock()

		// lk3 blocks.
		lk3ch := make(chan bool)
		go func() {
			lk3ch <- nsLk.lock(ctx, "volume", "path", "source", "opsID", false, 0)
		}()

		// lk4, blocks.
		lk4ch := make(chan bool)
		go func() {
			lk4ch <- nsLk.lock(ctx, "volume", "path", "source", "opsID", false, 0)
		}()
		runtime.Gosched()

		// unlock the manual lock
		nsLk.lockMapMutex.Unlock()

		// To trigger the race:
		// 1) lk3 or lk4 need to advance and increment the ref on the existing resource,
		//    successfully acquiring the lock.
		// 2) lk2 then needs to advance and remove the resource from lockMap.
		// 3) lk3 or lk4 (whichever didn't execute in step 1) then executes and creates
		//    a new entry in lockMap and acquires a lock for the same resource.

		<-lk2ch
		lk3ok := <-lk3ch
		lk4ok := <-lk4ch

		if lk3ok && lk4ok {
			t.Fatalf("multiple locks acquired; iteration=%d, lk3=%t, lk4=%t", i, lk3ok, lk4ok)
		}
	}
}
