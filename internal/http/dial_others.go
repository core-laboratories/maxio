//go:build !linux
// +build !linux

// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package http

import (
	"context"
	"net"
	"syscall"
	"time"
)

// TODO: if possible implement for non-linux platforms, not a priority at the moment
//
//nolint:unused
func setTCPParametersFn(opts TCPOptions) func(network, address string, c syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		return nil
	}
}

// DialContext is a function to make custom Dial for internode communications
type DialContext func(ctx context.Context, network, address string) (net.Conn, error)

// NewInternodeDialContext configures a custom dialer for internode communications
func NewInternodeDialContext(dialTimeout time.Duration, _ TCPOptions) DialContext {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout: dialTimeout,
		}
		return dialer.DialContext(ctx, network, addr)
	}
}
