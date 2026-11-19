// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package opa

import (
	"github.com/core-laboratories/maxio/internal/config"
)

// Legacy OPA envs
const (
	EnvIamOpaURL       = "MINIO_IAM_OPA_URL"
	EnvIamOpaAuthToken = "MINIO_IAM_OPA_AUTHTOKEN"
)

// SetPolicyOPAConfig - One time migration code needed, for migrating from older config to new for PolicyOPAConfig.
func SetPolicyOPAConfig(s config.Config, opaArgs Args) {
	if opaArgs.URL == nil || opaArgs.URL.String() == "" {
		// Do not enable if opaArgs was empty.
		return
	}
	s[config.PolicyOPASubSys][config.Default] = config.KVS{
		config.KV{
			Key:   URL,
			Value: opaArgs.URL.String(),
		},
		config.KV{
			Key:   AuthToken,
			Value: opaArgs.AuthToken,
		},
	}
}
