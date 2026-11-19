// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the administrator-controlled IPFS integration configuration.
// IPFS is enabled unless Enabled is set to false explicitly.
type Config struct {
	Enabled         bool
	AutoPin         bool
	ClusterEndpoint string
	TenantBucket    string
	ClusterUser     string
	ClusterPassword string
	KuboEndpoint    string
	KuboToken       string
	Replication     int
	Chunker         string
	RequestTimeout  time.Duration
	DefaultAccess   Access
	StatusHeaders   bool
	PublicGateway   bool
	GatewayListen   string
	GatewayDomain   string
	GatewayRate     int64
}

func ConfigFromEnv() (Config, error) {
	c := Config{
		ClusterEndpoint: os.Getenv("MAXIO_IPFS_CLUSTER_ENDPOINT"),
		TenantBucket:    os.Getenv("MAXIO_IPFS_TENANT_BUCKET"),
		ClusterUser:     os.Getenv("MAXIO_IPFS_CLUSTER_USER"),
		ClusterPassword: os.Getenv("MAXIO_IPFS_CLUSTER_PASSWORD"),
		KuboEndpoint:    os.Getenv("MAXIO_IPFS_KUBO_ENDPOINT"),
		KuboToken:       os.Getenv("MAXIO_IPFS_KUBO_TOKEN"),
		Replication:     2,
		Chunker:         "size-262144",
		RequestTimeout:  30 * time.Minute,
		DefaultAccess:   Private,
		GatewayListen:   "127.0.0.1:8090",
		GatewayRate:     10 << 20,
	}
	var err error
	if c.Enabled, err = envBool("MAXIO_IPFS_ENABLED", true); err != nil {
		return Config{}, err
	}
	if c.AutoPin, err = envBool("MAXIO_IPFS_AUTO_PIN", false); err != nil {
		return Config{}, err
	}
	if c.StatusHeaders, err = envBool("MAXIO_IPFS_STATUS_HEADERS", true); err != nil {
		return Config{}, err
	}
	if c.PublicGateway, err = envBool("MAXIO_IPFS_PUBLIC_GATEWAY", false); err != nil {
		return Config{}, err
	}
	if v := os.Getenv("MAXIO_IPFS_GATEWAY_LISTEN"); v != "" {
		c.GatewayListen = v
	}
	if v := os.Getenv("MAXIO_IPFS_GATEWAY_DOMAIN"); v != "" {
		c.GatewayDomain = strings.ToLower(v)
	}
	if v := os.Getenv("MAXIO_IPFS_GATEWAY_BYTES_PER_SECOND"); v != "" {
		c.GatewayRate, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("MAXIO_IPFS_GATEWAY_BYTES_PER_SECOND: %w", err)
		}
	}
	if v := os.Getenv("MAXIO_IPFS_REPLICATION"); v != "" {
		c.Replication, err = strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("MAXIO_IPFS_REPLICATION: %w", err)
		}
	}
	if v := os.Getenv("MAXIO_IPFS_IMPORT_CHUNKER"); v != "" {
		c.Chunker = v
	}
	if v := os.Getenv("MAXIO_IPFS_DEFAULT_ACCESS"); v != "" {
		c.DefaultAccess = Access(v)
	}
	return c, c.Validate()
}

func envBool(name string, fallback bool) (bool, error) {
	v := os.Getenv(name)
	if v == "" {
		return fallback, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("%s: %w", name, err)
	}
	return b, nil
}

func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Replication != 2 {
		return errors.New("IPFS currently requires replication factor 2")
	}
	if c.TenantBucket == "" {
		return errors.New("MAXIO_IPFS_TENANT_BUCKET is required")
	}
	if c.DefaultAccess != Private && c.DefaultAccess != Public {
		return errors.New("invalid IPFS default access")
	}
	if c.AutoPin && c.DefaultAccess == Private {
		return errors.New("automatic pinning of private plaintext is unsupported; set default access public or disable automatic pinning")
	}
	if c.Chunker != "size-262144" {
		return errors.New("unsupported IPFS chunker; changing it would change existing CIDs")
	}
	if c.PublicGateway && (c.GatewayListen == "" || c.GatewayRate <= 0) {
		return errors.New("public IPFS gateway requires a listen address and positive bandwidth limit")
	}
	if c.PublicGateway {
		host, _, err := net.SplitHostPort(c.GatewayListen)
		if err != nil {
			return fmt.Errorf("gateway listen address: %w", err)
		}
		if !strings.EqualFold(host, "localhost") && (net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback()) {
			return errors.New("public IPFS gateway must listen on loopback behind a TLS reverse proxy")
		}
	}
	if err := validateInternalEndpoint(c.ClusterEndpoint); err != nil {
		return fmt.Errorf("cluster endpoint: %w", err)
	}
	if err := validateInternalEndpoint(c.KuboEndpoint); err != nil {
		return fmt.Errorf("Kubo endpoint: %w", err)
	}
	if u, _ := url.Parse(c.KuboEndpoint); u != nil && c.KuboToken == "" {
		host := u.Hostname()
		ip := net.ParseIP(host)
		if !strings.EqualFold(host, "localhost") && (ip == nil || !ip.IsLoopback()) {
			return errors.New("remote Kubo RPC requires a bearer token")
		}
	}
	if c.ClusterUser == "" || c.ClusterPassword == "" {
		return errors.New("cluster basic authentication is required")
	}
	return nil
}

func validateInternalEndpoint(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("expected an absolute endpoint URL without credentials or query")
	}
	if u.Scheme == "https" {
		return nil
	}
	if u.Scheme != "http" {
		return errors.New("endpoint must use HTTPS or loopback HTTP")
	}
	host := u.Hostname()
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return nil
	}
	return errors.New("non-loopback IPFS RPC requires HTTPS")
}
