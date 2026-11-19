// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"strconv"

	"github.com/core-laboratories/maxio/internal/bucket/lifecycle"
)

//go:generate stringer -type lcEventSrc -trimprefix lcEventSrc_ $GOFILE
type lcEventSrc uint8

//nolint:staticcheck,revive // Underscores are used here to indicate where common prefix ends and the enumeration name begins
const (
	lcEventSrc_None lcEventSrc = iota
	lcEventSrc_Heal
	lcEventSrc_Scanner
	lcEventSrc_Decom
	lcEventSrc_Rebal
	lcEventSrc_s3HeadObject
	lcEventSrc_s3GetObject
	lcEventSrc_s3ListObjects
	lcEventSrc_s3PutObject
	lcEventSrc_s3CopyObject
	lcEventSrc_s3CompleteMultipartUpload
)

//revive:enable:var-naming
type lcAuditEvent struct {
	lifecycle.Event
	source lcEventSrc
}

func (lae lcAuditEvent) Tags() map[string]string {
	event := lae.Event
	src := lae.source
	const (
		ilmSrc                     = "ilm-src"
		ilmAction                  = "ilm-action"
		ilmDue                     = "ilm-due"
		ilmRuleID                  = "ilm-rule-id"
		ilmTier                    = "ilm-tier"
		ilmNewerNoncurrentVersions = "ilm-newer-noncurrent-versions"
		ilmNoncurrentDays          = "ilm-noncurrent-days"
	)
	tags := make(map[string]string, 5)
	if src > lcEventSrc_None {
		tags[ilmSrc] = src.String()
	}
	tags[ilmAction] = event.Action.String()
	tags[ilmRuleID] = event.RuleID

	if !event.Due.IsZero() {
		tags[ilmDue] = event.Due.Format(iso8601Format)
	}

	// rule with Transition/NoncurrentVersionTransition in effect
	if event.StorageClass != "" {
		tags[ilmTier] = event.StorageClass
	}

	// rule with NewernoncurrentVersions in effect
	if event.NewerNoncurrentVersions > 0 {
		tags[ilmNewerNoncurrentVersions] = strconv.Itoa(event.NewerNoncurrentVersions)
	}
	if event.NoncurrentDays > 0 {
		tags[ilmNoncurrentDays] = strconv.Itoa(event.NoncurrentDays)
	}
	return tags
}

func newLifecycleAuditEvent(src lcEventSrc, event lifecycle.Event) lcAuditEvent {
	return lcAuditEvent{
		Event:  event,
		source: src,
	}
}
