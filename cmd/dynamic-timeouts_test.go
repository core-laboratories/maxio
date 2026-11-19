// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestDynamicTimeoutSingleIncrease(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	initial := timeout.Timeout()

	for range dynamicTimeoutLogSize {
		timeout.LogFailure()
	}

	adjusted := timeout.Timeout()

	if initial >= adjusted {
		t.Errorf("Failure to increase timeout, expected %v to be more than %v", adjusted, initial)
	}
}

func TestDynamicTimeoutDualIncrease(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	initial := timeout.Timeout()

	for range dynamicTimeoutLogSize {
		timeout.LogFailure()
	}

	adjusted := timeout.Timeout()

	for range dynamicTimeoutLogSize {
		timeout.LogFailure()
	}

	adjustedAgain := timeout.Timeout()

	if initial >= adjusted || adjusted >= adjustedAgain {
		t.Errorf("Failure to increase timeout multiple times")
	}
}

func TestDynamicTimeoutSingleDecrease(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	initial := timeout.Timeout()

	for range dynamicTimeoutLogSize {
		timeout.LogSuccess(20 * time.Second)
	}

	adjusted := timeout.Timeout()

	if initial <= adjusted {
		t.Errorf("Failure to decrease timeout, expected %v to be less than %v", adjusted, initial)
	}
}

func TestDynamicTimeoutDualDecrease(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	initial := timeout.Timeout()

	for range dynamicTimeoutLogSize {
		timeout.LogSuccess(20 * time.Second)
	}

	adjusted := timeout.Timeout()

	for range dynamicTimeoutLogSize {
		timeout.LogSuccess(20 * time.Second)
	}

	adjustedAgain := timeout.Timeout()

	if initial <= adjusted || adjusted <= adjustedAgain {
		t.Errorf("Failure to decrease timeout multiple times, initial: %v, adjusted: %v, again: %v", initial, adjusted, adjustedAgain)
	}
}

func TestDynamicTimeoutManyDecreases(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	initial := timeout.Timeout()

	const successTimeout = 20 * time.Second
	for range 100 {
		for range dynamicTimeoutLogSize {
			timeout.LogSuccess(successTimeout)
		}
	}

	adjusted := timeout.Timeout()
	// Check whether eventual timeout is between initial value and success timeout
	if initial <= adjusted || adjusted <= successTimeout {
		t.Errorf("Failure to decrease timeout appropriately")
	}
}

func TestDynamicTimeoutConcurrent(t *testing.T) {
	// Race test.
	timeout := newDynamicTimeout(time.Second, time.Millisecond)
	var wg sync.WaitGroup
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		rng := rand.New(rand.NewSource(int64(i)))
		go func() {
			defer wg.Done()
			for range 100 {
				for range 100 {
					timeout.LogSuccess(time.Duration(float64(time.Second) * rng.Float64()))
				}
				to := timeout.Timeout()
				if to < time.Millisecond || to > time.Second {
					panic(to)
				}
			}
		}()
	}
	wg.Wait()
}

func TestDynamicTimeoutHitMinimum(t *testing.T) {
	const minimum = 30 * time.Second
	timeout := newDynamicTimeout(time.Minute, minimum)

	initial := timeout.Timeout()

	const successTimeout = 20 * time.Second
	for range 100 {
		for range dynamicTimeoutLogSize {
			timeout.LogSuccess(successTimeout)
		}
	}

	adjusted := timeout.Timeout()
	// Check whether eventual timeout has hit the minimum value
	if initial <= adjusted || adjusted != minimum {
		t.Errorf("Failure to decrease timeout appropriately")
	}
}

func testDynamicTimeoutAdjust(t *testing.T, timeout *dynamicTimeout, f func() float64) {
	const successTimeout = 20 * time.Second

	for range dynamicTimeoutLogSize {
		rnd := f()
		duration := max(time.Duration(float64(successTimeout)*rnd), 100*time.Millisecond)
		if duration >= time.Minute {
			timeout.LogFailure()
		} else {
			timeout.LogSuccess(duration)
		}
	}
}

func TestDynamicTimeoutAdjustExponential(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	rand.Seed(0)

	initial := timeout.Timeout()

	for range 10 {
		testDynamicTimeoutAdjust(t, timeout, rand.ExpFloat64)
	}

	adjusted := timeout.Timeout()
	if initial <= adjusted {
		t.Errorf("Failure to decrease timeout, expected %v to be less than %v", adjusted, initial)
	}
}

func TestDynamicTimeoutAdjustNormalized(t *testing.T) {
	timeout := newDynamicTimeout(time.Minute, time.Second)

	rand.Seed(0)

	initial := timeout.Timeout()

	for range 10 {
		testDynamicTimeoutAdjust(t, timeout, func() float64 {
			return 1.0 + rand.NormFloat64()
		})
	}

	adjusted := timeout.Timeout()
	if initial <= adjusted {
		t.Errorf("Failure to decrease timeout, expected %v to be less than %v", adjusted, initial)
	}
}
