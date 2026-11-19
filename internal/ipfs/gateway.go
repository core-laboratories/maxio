// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PublicGateway serves only verified, public pins backed by authorized S3 objects.
// It never forwards arbitrary CID requests to Kubo.
type PublicGateway struct {
	Manager *Manager
	Slots   chan struct{}
}

func NewPublicGateway(manager *Manager) *PublicGateway {
	return &PublicGateway{Manager: manager, Slots: make(chan struct{}, 16)}
}

func (g *PublicGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	select {
	case g.Slots <- struct{}{}:
		defer func() { <-g.Slots }()
	default:
		http.Error(w, "busy", http.StatusTooManyRequests)
		return
	}
	cid := strings.TrimPrefix(r.URL.Path, "/ipfs/")
	if cid == r.URL.Path && g.Manager.Config.GatewayDomain != "" {
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		suffix := ".ipfs." + g.Manager.Config.GatewayDomain
		if strings.HasSuffix(strings.ToLower(host), suffix) {
			cid = strings.TrimSuffix(strings.ToLower(host), suffix)
		}
	}
	if !cidPattern.MatchString(cid) {
		http.NotFound(w, r)
		return
	}
	jobs, err := g.Manager.Store.List(r.Context())
	if err != nil {
		http.Error(w, "gateway unavailable", http.StatusServiceUnavailable)
		return
	}
	var chosen *Job
	var reader io.ReadCloser
	for i := range jobs {
		job := &jobs[i]
		if job.CID == cid && job.DesiredPin && job.Access == Public && job.Status == Pinned && job.Replicas >= g.Manager.Config.Replication {
			candidate, err := g.Manager.Source.Open(r.Context(), job.Object)
			if err == nil {
				chosen = job
				reader = candidate
				break
			}
		}
	}
	if chosen == nil {
		http.NotFound(w, r)
		return
	}
	defer reader.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="content"`)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "sandbox")
	w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("ETag", `"`+cid+`"`)
	start, length, partial, err := gatewayRange(r.Header.Get("Range"), chosen.Object.Size)
	if err != nil {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", chosen.Object.Size))
		http.Error(w, "invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}
	w.Header().Set("Accept-Ranges", "bytes")
	if length >= 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(length, 10))
	}
	if partial {
		if _, err := io.CopyN(io.Discard, reader, start); err != nil {
			http.Error(w, "content unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, start+length-1, chosen.Object.Size))
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if r.Method == http.MethodHead {
		return
	}
	writer := &bandwidthWriter{ctx: r.Context(), writer: w, limit: g.Manager.Config.GatewayRate, begin: time.Now()}
	if length >= 0 {
		_, _ = io.CopyN(writer, reader, length)
	} else {
		_, _ = io.Copy(writer, reader)
	}
}

func gatewayRange(raw string, size int64) (start, length int64, partial bool, err error) {
	if raw == "" {
		return 0, size, false, nil
	}
	if size <= 0 || !strings.HasPrefix(raw, "bytes=") || strings.Contains(raw, ",") {
		return 0, 0, false, fmt.Errorf("invalid range")
	}
	parts := strings.SplitN(strings.TrimPrefix(raw, "bytes="), "-", 2)
	if len(parts) != 2 {
		return 0, 0, false, fmt.Errorf("invalid range")
	}
	if parts[0] == "" {
		last, e := strconv.ParseInt(parts[1], 10, 64)
		if e != nil || last <= 0 {
			return 0, 0, false, fmt.Errorf("invalid suffix")
		}
		if last > size {
			last = size
		}
		return size - last, last, true, nil
	}
	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 || start >= size {
		return 0, 0, false, fmt.Errorf("invalid start")
	}
	end := size - 1
	if parts[1] != "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end < start {
			return 0, 0, false, fmt.Errorf("invalid end")
		}
		if end >= size {
			end = size - 1
		}
	}
	return start, end - start + 1, true, nil
}

type bandwidthWriter struct {
	ctx     context.Context
	writer  io.Writer
	limit   int64
	written int64
	begin   time.Time
}

func (w *bandwidthWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	w.written += int64(n)
	if w.limit > 0 {
		deadline := w.begin.Add(time.Duration(float64(w.written) / float64(w.limit) * float64(time.Second)))
		if delay := time.Until(deadline); delay > 0 {
			select {
			case <-time.After(delay):
			case <-w.ctx.Done():
				return n, w.ctx.Err()
			}
		}
	}
	return n, err
}
