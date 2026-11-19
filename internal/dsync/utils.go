// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package dsync

import (
	"math/rand"
	"time"
)

func backoffWait(minSleep, unit, maxSleep time.Duration) func(*rand.Rand, uint) time.Duration {
	if unit > time.Hour {
		// Protect against integer overflow
		panic("unit cannot exceed one hour")
	}
	return func(r *rand.Rand, attempt uint) time.Duration {
		sleep := minSleep
		sleep += unit * time.Duration(attempt)
		if sleep > maxSleep {
			sleep = maxSleep
		}
		sleep -= time.Duration(r.Float64() * float64(sleep))
		return sleep
	}
}
