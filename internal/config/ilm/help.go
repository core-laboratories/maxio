// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ilm

import "github.com/core-laboratories/maxio/internal/config"

const (
	transitionWorkers = "transition_workers"
	expirationWorkers = "expiration_workers"
	// EnvILMTransitionWorkers env variable to configure number of transition workers
	EnvILMTransitionWorkers = "MINIO_ILM_TRANSITION_WORKERS"
	// EnvILMExpirationWorkers env variable to configure number of expiration workers
	EnvILMExpirationWorkers = "MINIO_ILM_EXPIRATION_WORKERS"
)

var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	// Help holds configuration keys and their default values for the ILM
	// subsystem
	Help = config.HelpKVS{
		config.HelpKV{
			Key:         transitionWorkers,
			Type:        "number",
			Description: `set the number of transition workers` + defaultHelpPostfix(transitionWorkers),
			Optional:    true,
		},
		config.HelpKV{
			Key:         expirationWorkers,
			Type:        "number",
			Description: `set the number of expiration workers` + defaultHelpPostfix(expirationWorkers),
			Optional:    true,
		},
	}
)
