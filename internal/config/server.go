// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package config

// Opts holds MinIO configuration options
type Opts struct {
	FTP struct {
		Address          string `yaml:"address"`
		PassivePortRange string `yaml:"passive-port-range"`
	} `yaml:"ftp"`
	SFTP struct {
		Address       string `yaml:"address"`
		SSHPrivateKey string `yaml:"ssh-private-key"`
	} `yaml:"sftp"`
}

// ServerConfigVersion struct is used to extract the version
type ServerConfigVersion struct {
	Version string `yaml:"version"`
}

// ServerConfigCommon struct for server config common options
type ServerConfigCommon struct {
	RootUser    string `yaml:"rootUser"`
	RootPwd     string `yaml:"rootPassword"`
	Addr        string `yaml:"address"`
	ConsoleAddr string `yaml:"console-address"`
	CertsDir    string `yaml:"certs-dir"`
	Options     Opts   `yaml:"options"`
}

// ServerConfigV1 represents a MinIO configuration file v1
type ServerConfigV1 struct {
	ServerConfigVersion
	ServerConfigCommon
	Pools [][]string `yaml:"pools"`
}

// ServerConfig represents a MinIO configuration file
type ServerConfig struct {
	ServerConfigVersion
	ServerConfigCommon
	Pools []struct {
		Args          []string `yaml:"args"`
		SetDriveCount uint64   `yaml:"set-drive-count"`
	} `yaml:"pools"`
}
