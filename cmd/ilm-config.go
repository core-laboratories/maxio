// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"sync"

	"github.com/core-laboratories/maxio/internal/config/ilm"
)

var globalILMConfig = ilmConfig{
	cfg: ilm.Config{
		ExpirationWorkers: 100,
		TransitionWorkers: 100,
	},
}

type ilmConfig struct {
	mu  sync.RWMutex
	cfg ilm.Config
}

func (c *ilmConfig) getExpirationWorkers() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cfg.ExpirationWorkers
}

func (c *ilmConfig) getTransitionWorkers() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cfg.TransitionWorkers
}

func (c *ilmConfig) update(cfg ilm.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cfg = cfg
}
