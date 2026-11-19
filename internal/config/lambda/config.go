// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package lambda

import "github.com/core-laboratories/maxio/internal/event/target"

// Config - lambda target configuration structure, holds
// information about various lambda targets.
type Config struct {
	Webhook map[string]target.WebhookArgs `json:"webhook"`
}

const (
	defaultTarget = "1"
)

// NewConfig - initialize lambda config.
func NewConfig() Config {
	// Make sure to initialize lambda targets
	cfg := Config{
		Webhook: make(map[string]target.WebhookArgs),
	}
	cfg.Webhook[defaultTarget] = target.WebhookArgs{}
	return cfg
}
