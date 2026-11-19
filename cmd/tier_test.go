// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"testing"
	"time"
)

func TestTierMetrics(t *testing.T) {
	tier := "WARM-1"
	globalTierMetrics.Observe(tier, 200*time.Millisecond)
	expSuccess := 10
	expFailure := 5
	for range expSuccess {
		globalTierMetrics.logSuccess(tier)
	}
	for range expFailure {
		globalTierMetrics.logFailure(tier)
	}
	metrics := globalTierMetrics.Report()
	var succ, fail float64
	for _, metric := range metrics {
		switch metric.Description.Name {
		case tierRequestsSuccess:
			succ += metric.Value
		case tierRequestsFailure:
			fail += metric.Value
		}
	}
	if int(succ) != expSuccess {
		t.Fatalf("Expected %d successes but got %f", expSuccess, succ)
	}
	if int(fail) != expFailure {
		t.Fatalf("Expected %d failures but got %f", expFailure, fail)
	}
}
