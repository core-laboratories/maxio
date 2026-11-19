// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package bpool

import "sync"

// Pool is a single type sync.Pool with a few extra properties:
// If New is not set Get may return the zero value of T.
type Pool[T any] struct {
	New func() T
	p   sync.Pool
}

// Get will retuen a new T
func (p *Pool[T]) Get() T {
	v, ok := p.p.Get().(T)
	if ok {
		return v
	}
	if p.New == nil {
		var t T
		return t
	}
	return p.New()
}

// Put a used T.
func (p *Pool[T]) Put(t T) {
	p.p.Put(t)
}
