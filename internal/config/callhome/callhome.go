// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package callhome

import (
	"sync"
	"time"

	"github.com/core-laboratories/maxio/internal/config"
	"github.com/minio/pkg/v3/env"
)

// Callhome related keys
const (
	Enable    = "enable"
	Frequency = "frequency"
)

// DefaultKVS - default KV config for subnet settings
var DefaultKVS = config.KVS{
	config.KV{
		Key:   Enable,
		Value: "off",
	},
	config.KV{
		Key:   Frequency,
		Value: "24h",
	},
}

// callhomeCycleDefault is the default interval between two callhome cycles (24hrs)
const callhomeCycleDefault = 24 * time.Hour

// Config represents the subnet related configuration
type Config struct {
	// Flag indicating whether callhome is enabled.
	Enable bool `json:"enable"`

	// The interval between callhome cycles
	Frequency time.Duration `json:"frequency"`
}

var configLock sync.RWMutex

// Enabled - indicates if callhome is enabled or not
func (c *Config) Enabled() bool {
	configLock.RLock()
	defer configLock.RUnlock()

	return c.Enable
}

// FrequencyDur - returns the currently configured callhome frequency
func (c *Config) FrequencyDur() time.Duration {
	configLock.RLock()
	defer configLock.RUnlock()

	if c.Frequency == 0 {
		return callhomeCycleDefault
	}

	return c.Frequency
}

// Update updates new callhome frequency
func (c *Config) Update(ncfg Config) {
	configLock.Lock()
	defer configLock.Unlock()

	c.Enable = ncfg.Enable
	c.Frequency = ncfg.Frequency
}

// LookupConfig - lookup config and override with valid environment settings if any.
func LookupConfig(kvs config.KVS) (cfg Config, err error) {
	if err = config.CheckValidKeys(config.CallhomeSubSys, kvs, DefaultKVS); err != nil {
		return cfg, err
	}

	cfg.Enable = env.Get(config.EnvMinIOCallhomeEnable,
		kvs.GetWithDefault(Enable, DefaultKVS)) == config.EnableOn
	cfg.Frequency, err = time.ParseDuration(env.Get(config.EnvMinIOCallhomeFrequency,
		kvs.GetWithDefault(Frequency, DefaultKVS)))
	return cfg, err
}
