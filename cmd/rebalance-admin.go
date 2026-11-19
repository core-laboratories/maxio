// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"time"
)

type rebalPoolProgress struct {
	NumObjects  uint64        `json:"objects"`
	NumVersions uint64        `json:"versions"`
	Bytes       uint64        `json:"bytes"`
	Bucket      string        `json:"bucket"`
	Object      string        `json:"object"`
	Elapsed     time.Duration `json:"elapsed"`
	ETA         time.Duration `json:"eta"`
}

type rebalancePoolStatus struct {
	ID       int               `json:"id"`       // Pool index (zero-based)
	Status   string            `json:"status"`   // Active if rebalance is running, empty otherwise
	Used     float64           `json:"used"`     // Percentage used space
	Progress rebalPoolProgress `json:"progress"` // is empty when rebalance is not running
}

// rebalanceAdminStatus holds rebalance status related information exported to mc, console, etc.
type rebalanceAdminStatus struct {
	ID        string                // identifies the ongoing rebalance operation by a uuid
	Pools     []rebalancePoolStatus `json:"pools"` // contains all pools, including inactive
	StoppedAt time.Time             `json:"stoppedAt"`
}

func rebalanceStatus(ctx context.Context, z *erasureServerPools) (r rebalanceAdminStatus, err error) {
	// Load latest rebalance status
	meta := &rebalanceMeta{}
	err = meta.load(ctx, z.serverPools[0])
	if err != nil {
		return r, err
	}

	// Compute disk usage percentage
	si := z.StorageInfo(ctx, true)
	diskStats := make([]struct {
		AvailableSpace uint64
		TotalSpace     uint64
	}, len(z.serverPools))
	for _, disk := range si.Disks {
		// Ignore invalid.
		if disk.PoolIndex < 0 || len(diskStats) <= disk.PoolIndex {
			// https://github.com/core-laboratories/maxio/issues/16500
			continue
		}
		diskStats[disk.PoolIndex].AvailableSpace += disk.AvailableSpace
		diskStats[disk.PoolIndex].TotalSpace += disk.TotalSpace
	}

	stopTime := meta.StoppedAt
	r = rebalanceAdminStatus{
		ID:        meta.ID,
		StoppedAt: meta.StoppedAt,
		Pools:     make([]rebalancePoolStatus, len(meta.PoolStats)),
	}
	for i, ps := range meta.PoolStats {
		r.Pools[i] = rebalancePoolStatus{
			ID:     i,
			Status: ps.Info.Status.String(),
			Used:   float64(diskStats[i].TotalSpace-diskStats[i].AvailableSpace) / float64(diskStats[i].TotalSpace),
		}
		if !ps.Participating {
			continue
		}
		// for participating pools, total bytes to be rebalanced by this pool is given by,
		// pf_c = (f_i + x)/c_i,
		// pf_c - percentage free space across pools, f_i - ith pool's free space, c_i - ith pool's capacity
		// i.e. x = c_i*pfc -f_i
		totalBytesToRebal := float64(ps.InitCapacity)*meta.PercentFreeGoal - float64(ps.InitFreeSpace)
		elapsed := time.Since(ps.Info.StartTime)
		eta := time.Duration(totalBytesToRebal * float64(elapsed) / float64(ps.Bytes))
		if !ps.Info.EndTime.IsZero() {
			stopTime = ps.Info.EndTime
		}

		if !stopTime.IsZero() { // rebalance is stopped or completed
			elapsed = stopTime.Sub(ps.Info.StartTime)
			eta = 0
		}

		r.Pools[i].Progress = rebalPoolProgress{
			NumObjects:  ps.NumObjects,
			NumVersions: ps.NumVersions,
			Bytes:       ps.Bytes,
			Elapsed:     elapsed,
			ETA:         eta,
		}
	}
	return r, nil
}
