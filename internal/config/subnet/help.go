// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package subnet

import "github.com/core-laboratories/maxio/internal/config"

var (
	defaultHelpPostfix = func(key string) string {
		return config.DefaultHelpPostfix(DefaultKVS, key)
	}

	// HelpSubnet - provides help for subnet api key config
	HelpSubnet = config.HelpKVS{
		config.HelpKV{
			Key:         config.License,
			Type:        "string",
			Description: "Enterprise license for the cluster" + defaultHelpPostfix(config.License),
			Optional:    true,
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         config.APIKey,
			Type:        "string",
			Description: "Enterprise license API key for the cluster" + defaultHelpPostfix(config.APIKey),
			Optional:    true,
			Sensitive:   true,
		},
		config.HelpKV{
			Key:         config.Proxy,
			Type:        "string",
			Description: "HTTP(s) proxy URL to use for connecting to SUBNET" + defaultHelpPostfix(config.Proxy),
			Optional:    true,
			Sensitive:   true,
		},
	}
)
