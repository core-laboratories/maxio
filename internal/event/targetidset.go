// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package event

import "maps"

// TargetIDSet - Set representation of TargetIDs.
type TargetIDSet map[TargetID]struct{}

// Clone - returns copy of this set.
func (set TargetIDSet) Clone() TargetIDSet {
	setCopy := NewTargetIDSet()
	maps.Copy(setCopy, set)
	return setCopy
}

// add - adds TargetID to the set.
func (set TargetIDSet) add(targetID TargetID) {
	set[targetID] = struct{}{}
}

// Union - returns union with given set as new set.
func (set TargetIDSet) Union(sset TargetIDSet) TargetIDSet {
	nset := set.Clone()

	for k := range sset {
		nset.add(k)
	}

	return nset
}

// Difference - returns difference with given set as new set.
func (set TargetIDSet) Difference(sset TargetIDSet) TargetIDSet {
	nset := NewTargetIDSet()
	for k := range set {
		if _, ok := sset[k]; !ok {
			nset.add(k)
		}
	}

	return nset
}

// NewTargetIDSet - creates new TargetID set with given TargetIDs.
func NewTargetIDSet(targetIDs ...TargetID) TargetIDSet {
	set := make(TargetIDSet)
	for _, targetID := range targetIDs {
		set.add(targetID)
	}
	return set
}
