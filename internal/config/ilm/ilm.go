// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ilm

import (
	"strconv"

	"github.com/core-laboratories/maxio/internal/config"
	"github.com/minio/pkg/v3/env"
)

// DefaultKVS default configuration values for ILM subsystem
var DefaultKVS = config.KVS{
	config.KV{
		Key:   transitionWorkers,
		Value: "100",
	},
	config.KV{
		Key:   expirationWorkers,
		Value: "100",
	},
}

// Config represents the different configuration values for ILM subsystem
type Config struct {
	TransitionWorkers int
	ExpirationWorkers int
}

// LookupConfig - lookup ilm config and override with valid environment settings if any.
func LookupConfig(kvs config.KVS) (cfg Config, err error) {
	cfg = Config{
		TransitionWorkers: 100,
		ExpirationWorkers: 100,
	}

	if err = config.CheckValidKeys(config.ILMSubSys, kvs, DefaultKVS); err != nil {
		return cfg, err
	}

	tw, err := strconv.Atoi(env.Get(EnvILMTransitionWorkers, kvs.GetWithDefault(transitionWorkers, DefaultKVS)))
	if err != nil {
		return cfg, err
	}

	ew, err := strconv.Atoi(env.Get(EnvILMExpirationWorkers, kvs.GetWithDefault(expirationWorkers, DefaultKVS)))
	if err != nil {
		return cfg, err
	}

	cfg.TransitionWorkers = tw
	cfg.ExpirationWorkers = ew
	return cfg, nil
}
