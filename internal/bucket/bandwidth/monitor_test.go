// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package bandwidth

import (
	"reflect"
	"testing"
	"time"
)

const (
	oneMiB uint64 = 1024 * 1024
)

func TestMonitor_GetReport(t *testing.T) {
	type fields struct {
		activeBuckets map[BucketOptions]*bucketMeasurement
		endTime       time.Time
		update2       uint64
		endTime2      time.Time
	}
	start := time.Now()
	m0 := newBucketMeasurement(start)
	m0.incrementBytes(0)
	m1MiBPS := newBucketMeasurement(start)
	m1MiBPS.incrementBytes(oneMiB)

	test1Want := make(map[BucketOptions]Details)
	test1Want[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = Details{LimitInBytesPerSecond: 1024 * 1024, CurrentBandwidthInBytesPerSecond: 0}
	test1Want2 := make(map[BucketOptions]Details)
	test1Want2[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = Details{
		LimitInBytesPerSecond:            1024 * 1024,
		CurrentBandwidthInBytesPerSecond: (1024 * 1024) / start.Add(2*time.Second).Sub(start.Add(1*time.Second)).Seconds(),
	}

	test2Want := make(map[BucketOptions]Details)
	test2Want[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = Details{LimitInBytesPerSecond: 1024 * 1024, CurrentBandwidthInBytesPerSecond: float64(oneMiB)}
	test2Want2 := make(map[BucketOptions]Details)
	test2Want2[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = Details{
		LimitInBytesPerSecond:            1024 * 1024,
		CurrentBandwidthInBytesPerSecond: exponentialMovingAverage(betaBucket, float64(oneMiB), 2*float64(oneMiB)),
	}

	test1ActiveBuckets := make(map[BucketOptions]*bucketMeasurement)
	test1ActiveBuckets[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = m0
	test1ActiveBuckets2 := make(map[BucketOptions]*bucketMeasurement)
	test1ActiveBuckets2[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = m1MiBPS

	tests := []struct {
		name   string
		fields fields
		want   *BucketBandwidthReport
		want2  *BucketBandwidthReport
	}{
		{
			name: "ZeroToOne",
			fields: fields{
				activeBuckets: test1ActiveBuckets,
				endTime:       start.Add(1 * time.Second),
				update2:       oneMiB,
				endTime2:      start.Add(2 * time.Second),
			},
			want: &BucketBandwidthReport{
				BucketStats: test1Want,
			},
			want2: &BucketBandwidthReport{
				BucketStats: test1Want2,
			},
		},
		{
			name: "OneToTwo",
			fields: fields{
				activeBuckets: test1ActiveBuckets2,
				endTime:       start.Add(1 * time.Second),
				update2:       2 * oneMiB,
				endTime2:      start.Add(2 * time.Second),
			},
			want: &BucketBandwidthReport{
				BucketStats: test2Want,
			},
			want2: &BucketBandwidthReport{
				BucketStats: test2Want2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			thr := bucketThrottle{
				NodeBandwidthPerSec: 1024 * 1024,
			}
			th := make(map[BucketOptions]*bucketThrottle)
			th[BucketOptions{Name: "bucket", ReplicationARN: "arn"}] = &thr
			m := &Monitor{
				bucketsMeasurement: tt.fields.activeBuckets,
				bucketsThrottle:    th,
				NodeCount:          1,
			}
			m.bucketsMeasurement[BucketOptions{Name: "bucket", ReplicationARN: "arn"}].updateExponentialMovingAverage(tt.fields.endTime)
			got := m.GetReport(SelectBuckets())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReport() = %v, want %v", got, tt.want)
			}
			m.bucketsMeasurement[BucketOptions{Name: "bucket", ReplicationARN: "arn"}].incrementBytes(tt.fields.update2)
			m.bucketsMeasurement[BucketOptions{Name: "bucket", ReplicationARN: "arn"}].updateExponentialMovingAverage(tt.fields.endTime2)
			got = m.GetReport(SelectBuckets())
			if !reflect.DeepEqual(got.BucketStats, tt.want2.BucketStats) {
				t.Errorf("GetReport() = %v, want %v", got.BucketStats, tt.want2.BucketStats)
			}
		})
	}
}
