// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"github.com/core-laboratories/maxio/internal/auth"
	"github.com/core-laboratories/maxio/internal/config"
	"github.com/core-laboratories/maxio/internal/config/compress"
	xldap "github.com/core-laboratories/maxio/internal/config/identity/ldap"
	"github.com/core-laboratories/maxio/internal/config/identity/openid"
	"github.com/core-laboratories/maxio/internal/config/notify"
	"github.com/core-laboratories/maxio/internal/config/policy/opa"
	"github.com/core-laboratories/maxio/internal/config/storageclass"
	"github.com/core-laboratories/maxio/internal/logger"
	"github.com/minio/pkg/v3/quick"
)

// FileLogger is introduced to workaround the dependency about logrus
type FileLogger struct {
	Enable   bool   `json:"enable"`
	Filename string `json:"filename"`
}

// ConsoleLogger is introduced to workaround the dependency about logrus
type ConsoleLogger struct {
	Enable bool `json:"enable"`
}

// serverConfigV33 is just like version '32', removes clientID from NATS and MQTT, and adds queueDir, queueLimit in all notification targets.
type serverConfigV33 struct {
	quick.Config `json:"-"` // ignore interfaces

	Version string `json:"version"`

	// S3 API configuration.
	Credential auth.Credentials `json:"credential"`
	Region     string           `json:"region"`
	Worm       config.BoolFlag  `json:"worm"`

	// Storage class configuration
	StorageClass storageclass.Config `json:"storageclass"`

	// Notification queue configuration.
	Notify notify.Config `json:"notify"`

	// Logger configuration
	Logger logger.Config `json:"logger"`

	// Compression configuration
	Compression compress.Config `json:"compress"`

	// OpenID configuration
	OpenID openid.Config `json:"openid"`

	// External policy enforcements.
	Policy struct {
		// OPA configuration.
		OPA opa.Args `json:"opa"`

		// Add new external policy enforcements here.
	} `json:"policy"`

	LDAPServerConfig xldap.LegacyConfig `json:"ldapserverconfig"`
}
