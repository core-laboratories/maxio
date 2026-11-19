//go:build linux
// +build linux

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

// CheckPortAvailability - check if given host and port is already in use.
// Note: The check method tries to listen on given port and closes it.
// It is possible to have a disconnected client in this tiny window of time.
func CheckPortAvailability(host, port string, opts TCPOptions) (err error) {
	lc := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			c.Control(func(fdPtr uintptr) {
				if opts.Interface != "" {
					// When interface is specified look for specifically port availability on
					// the specified interface if any.
					_ = syscall.SetsockoptString(int(fdPtr), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, opts.Interface)
				}
			})
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	l, err := lc.Listen(ctx, "tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}

	// As we are able to listen on this network, the port is not in use.
	// Close the listener and continue check other networks.
	return l.Close()
}
