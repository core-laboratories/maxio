// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package grid

import (
	"errors"
	"fmt"
)

var (
	// ErrUnknownHandler is returned when an unknown handler is requested.
	ErrUnknownHandler = errors.New("unknown mux handler")

	// ErrHandlerAlreadyExists is returned when a handler is already registered.
	ErrHandlerAlreadyExists = errors.New("mux handler already exists")

	// ErrIncorrectSequence is returned when an out-of-sequence item is received.
	ErrIncorrectSequence = errors.New("out-of-sequence item received")
)

// ErrResponse is a remote error response.
type ErrResponse struct {
	msg string
}

func (e ErrResponse) Error() string {
	return fmt.Sprintf("remote: %s", e.msg)
}
