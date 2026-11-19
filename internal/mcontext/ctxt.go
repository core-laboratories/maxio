// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package mcontext

// Share a common context information between different
// packages in github.com/core-laboratories/maxio

import (
	xhttp "github.com/core-laboratories/maxio/internal/http"
)

// ContextTraceType represents the type of golang Context key
type ContextTraceType string

// ContextTraceKey is the key of TraceCtxt saved in a Golang context
const ContextTraceKey = ContextTraceType("ctx-trace-info")

// TraceCtxt holds related tracing data of a http request.
type TraceCtxt struct {
	RequestRecorder  *xhttp.RequestRecorder
	ResponseRecorder *xhttp.ResponseRecorder

	FuncName string
	AmzReqID string
}
