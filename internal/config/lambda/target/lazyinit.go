// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package target

import (
	"sync"
	"sync/atomic"
)

// Inspired from Golang sync.Once but it is only marked
// initialized when the provided function returns nil.

type lazyInit struct {
	done uint32
	m    sync.Mutex
}

func (l *lazyInit) Do(f func() error) error {
	if atomic.LoadUint32(&l.done) == 0 {
		return l.doSlow(f)
	}
	return nil
}

func (l *lazyInit) doSlow(f func() error) error {
	l.m.Lock()
	defer l.m.Unlock()
	if atomic.LoadUint32(&l.done) == 0 {
		if err := f(); err != nil {
			return err
		}
		// Mark as done only when f() is successful
		atomic.StoreUint32(&l.done, 1)
	}
	return nil
}
