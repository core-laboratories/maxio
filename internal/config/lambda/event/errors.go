// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import (
	"fmt"
)

// ErrUnknownRegion - unknown region error.
type ErrUnknownRegion struct {
	Region string
}

func (err ErrUnknownRegion) Error() string {
	return fmt.Sprintf("unknown region '%v'", err.Region)
}

// ErrARNNotFound - ARN not found error.
type ErrARNNotFound struct {
	ARN ARN
}

func (err ErrARNNotFound) Error() string {
	return fmt.Sprintf("ARN '%v' not found", err.ARN)
}

// ErrInvalidARN - invalid ARN error.
type ErrInvalidARN struct {
	ARN string
}

func (err ErrInvalidARN) Error() string {
	return fmt.Sprintf("invalid ARN '%v'", err.ARN)
}
