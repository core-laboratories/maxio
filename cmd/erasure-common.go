// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

func (er erasureObjects) getOnlineDisks() (newDisks []StorageAPI) {
	disks := er.getDisks()
	var wg sync.WaitGroup
	var mu sync.Mutex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, i := range r.Perm(len(disks)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if disks[i] == nil {
				return
			}
			di, err := disks[i].DiskInfo(context.Background(), DiskInfoOptions{})
			if err != nil || di.Healing {
				// - Do not consume disks which are not reachable
				//   unformatted or simply not accessible for some reason.
				//
				// - Do not consume disks which are being healed
				//
				// - Future: skip busy disks
				return
			}

			mu.Lock()
			newDisks = append(newDisks, disks[i])
			mu.Unlock()
		}()
	}
	wg.Wait()
	return newDisks
}

func (er erasureObjects) getOnlineLocalDisks() (newDisks []StorageAPI) {
	disks := er.getOnlineDisks()

	// Based on the random shuffling return back randomized disks.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, i := range r.Perm(len(disks)) {
		if disks[i] != nil && disks[i].IsLocal() {
			newDisks = append(newDisks, disks[i])
		}
	}

	return newDisks
}

func (er erasureObjects) getLocalDisks() (newDisks []StorageAPI) {
	disks := er.getDisks()
	// Based on the random shuffling return back randomized disks.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, i := range r.Perm(len(disks)) {
		if disks[i] != nil && disks[i].IsLocal() {
			newDisks = append(newDisks, disks[i])
		}
	}
	return newDisks
}
