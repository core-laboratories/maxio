// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"path"
	"reflect"
	"strings"
	"time"
)

type ObjectRef struct {
	Bucket   string    `json:"bucket"`
	Key      string    `json:"key"`
	Version  string    `json:"version"`
	Owner    string    `json:"owner"`
	ETag     string    `json:"etag,omitempty"`
	Size     int64     `json:"size,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

func (o ObjectRef) ID() string {
	h := sha256.New()
	for _, part := range []string{o.Bucket, o.Key, o.Version} {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func (o ObjectRef) Valid() bool {
	return o.Bucket != "" && o.Key != "" && o.Owner != "" && !strings.HasPrefix(o.Key, "../") && path.Clean(o.Key) != "."
}

type JobStatus string

const (
	Pending  JobStatus = "pending"
	Pinning  JobStatus = "pinning"
	Pinned   JobStatus = "pinned"
	Failed   JobStatus = "failed"
	Dead     JobStatus = "dead"
	Unpinned JobStatus = "unpinned"
)

type Job struct {
	Revision    uint64    `json:"revision"`
	Object      ObjectRef `json:"object"`
	DesiredPin  bool      `json:"desiredPin"`
	Access      Access    `json:"access"`
	BasePolicy  Policy    `json:"basePolicy"`
	CID         string    `json:"cid,omitempty"`
	PreviousCID string    `json:"previousCID,omitempty"`
	Status      JobStatus `json:"status"`
	Replicas    int       `json:"replicas"`
	Importer    string    `json:"importer,omitempty"`
	Attempts    int       `json:"attempts"`
	NextAttempt time.Time `json:"nextAttempt,omitempty"`
	Updated     time.Time `json:"updated"`
	Error       string    `json:"error,omitempty"`
}

type JobStore interface {
	Load(ctx context.Context, id string) (Job, bool, error)
	Save(ctx context.Context, job Job) error
	CompareAndSwap(ctx context.Context, expected uint64, job Job) (bool, error)
	List(ctx context.Context) ([]Job, error)
}

type ObjectSource interface {
	Open(ctx context.Context, object ObjectRef) (io.ReadCloser, error)
}

type Manager struct {
	Config   Config
	Store    JobStore
	Source   ObjectSource
	Importer Importer
	Pins     PinManager
	Now      func() time.Time
}

func NewManager(c Config, store JobStore, source ObjectSource) *Manager {
	return &Manager{Config: c, Store: store, Source: source, Importer: NewKuboClient(c), Pins: NewClusterClient(c), Now: time.Now}
}

// RecordCreated persists the intent before any network import is started.
func (m *Manager) RecordCreated(ctx context.Context, ref ObjectRef, policy Policy) error {
	return m.RecordPolicy(ctx, ref, policy, Policy{}, true)
}

func (m *Manager) RecordPolicy(ctx context.Context, ref ObjectRef, base, tags Policy, replaceBase bool) error {
	if !m.Config.Enabled {
		return nil
	}
	if !ref.Valid() {
		return errors.New("invalid IPFS object reference")
	}
	for attempt := 0; attempt < 5; attempt++ {
		job, found, err := m.Store.Load(ctx, ref.ID())
		if err != nil {
			return err
		}
		if !found {
			job = Job{Object: ref}
		}
		if found && job.Object.Owner != ref.Owner {
			return errors.New("IPFS object owner changed")
		}
		effectiveBase := base
		if !replaceBase && found {
			effectiveBase = job.BasePolicy
		}
		policy := effectiveBase
		if tags.Pin != nil {
			policy.Pin = tags.Pin
		}
		if tags.Access != nil {
			policy.Access = tags.Access
		}
		access := m.Config.DefaultAccess
		if policy.Access != nil {
			access = *policy.Access
		}
		desired := EffectivePin(m.Config.AutoPin, policy)
		if found && job.Object.ETag == ref.ETag && job.Object.Modified.Equal(ref.Modified) && job.DesiredPin == desired && job.Access == access && reflect.DeepEqual(job.BasePolicy, effectiveBase) && job.Status != Failed && job.Status != Dead {
			return nil
		}
		expected := job.Revision
		if job.CID != "" {
			job.PreviousCID = job.CID
		}
		job.CID = ""
		job.Replicas = 0
		job.Object = ref
		job.BasePolicy = effectiveBase
		job.DesiredPin = desired
		job.Access = access
		job.Status = Pending
		job.Attempts = 0
		job.Error = ""
		job.NextAttempt = time.Time{}
		job.Updated = m.Now()
		changed, err := m.Store.CompareAndSwap(ctx, expected, job)
		if err != nil {
			return err
		}
		if changed {
			return nil
		}
	}
	return errors.New("concurrent IPFS policy update")
}

// RecordRemoved must only be called for an object version actually removed by S3.
func (m *Manager) RecordRemoved(ctx context.Context, ref ObjectRef) error {
	if !m.Config.Enabled {
		return nil
	}
	for attempt := 0; attempt < 5; attempt++ {
		job, found, err := m.Store.Load(ctx, ref.ID())
		if err != nil || !found {
			return err
		}
		if ref.Owner != "" && job.Object.Owner != ref.Owner {
			return errors.New("IPFS object owner mismatch")
		}
		expected := job.Revision
		job.DesiredPin = false
		job.Status = Pending
		job.NextAttempt = time.Time{}
		job.Updated = m.Now()
		changed, err := m.Store.CompareAndSwap(ctx, expected, job)
		if err != nil {
			return err
		}
		if changed {
			return nil
		}
	}
	return errors.New("concurrent IPFS removal")
}

func (m *Manager) Status(ctx context.Context, ref ObjectRef) (Job, bool, error) {
	if !m.Config.Enabled {
		return Job{}, false, nil
	}
	job, found, err := m.Store.Load(ctx, ref.ID())
	if found && job.Object.Owner != ref.Owner {
		return Job{}, false, nil
	}
	return job, found, err
}

// RunOnce retries unfinished jobs; multiple callers must serialize this method.
func (m *Manager) RunOnce(ctx context.Context) error {
	if !m.Config.Enabled {
		return nil
	}
	jobs, err := m.Store.List(ctx)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if err := ctx.Err(); err != nil {
			return err
		}
		expected := job.Revision
		if job.Status == Dead {
			continue
		}
		if job.Status == Pinned && job.DesiredPin && job.Replicas >= m.Config.Replication && m.Now().Sub(job.Updated) < 5*time.Minute {
			continue
		}
		if job.Status == Unpinned && !job.DesiredPin {
			continue
		}
		if !job.NextAttempt.IsZero() && m.Now().Before(job.NextAttempt) {
			continue
		}
		if err := m.process(ctx, &job, jobs); err != nil {
			job.Attempts++
			job.Status = Failed
			job.Replicas = 0
			job.Error = err.Error()
			if job.Attempts >= 12 {
				job.Status = Dead
				job.NextAttempt = time.Time{}
			} else {
				seconds := math.Min(3600, math.Pow(2, float64(job.Attempts)))
				job.NextAttempt = m.Now().Add(time.Duration(seconds) * time.Second)
			}
		}
		job.Updated = m.Now()
		if _, err := m.Store.CompareAndSwap(ctx, expected, job); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) process(ctx context.Context, job *Job, all []Job) error {
	if !job.DesiredPin {
		waitForRefs := false
		for _, cid := range []string{job.CID, job.PreviousCID} {
			if cid == "" {
				continue
			}
			if referencedElsewhere(cid, job.Object.ID(), all) {
				waitForRefs = true
				continue
			}
			if err := m.Pins.Unpin(ctx, cid); err != nil {
				return err
			}
		}
		if waitForRefs {
			job.Status = Pending
			job.NextAttempt = m.Now().Add(5 * time.Second)
			return nil
		}
		job.Status = Unpinned
		job.Replicas = 0
		job.PreviousCID = ""
		return nil
	}
	if job.Access != Public {
		return errors.New("private content publishing requires an isolated or encrypted distribution mode")
	}
	if job.CID == "" {
		reader, err := m.Source.Open(ctx, job.Object)
		if err != nil {
			return err
		}
		job.Status = Pinning
		cid, err := m.Importer.Add(ctx, path.Base(job.Object.Key), reader)
		closeErr := reader.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
		verifyReader, err := m.Source.Open(ctx, job.Object)
		if err != nil {
			return err
		}
		expected, verifyErr := m.Importer.HashOnly(ctx, path.Base(job.Object.Key), verifyReader)
		closeErr = verifyReader.Close()
		if verifyErr != nil {
			return verifyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if cid != expected {
			return errors.New("IPFS CID verification failed")
		}
		job.CID = cid
		job.Importer = ImportVersion(m.Config)
	}
	if err := m.Pins.Pin(ctx, job.CID); err != nil {
		return err
	}
	replicas, err := m.Pins.Replicas(ctx, job.CID)
	if err != nil {
		return err
	}
	job.Replicas = replicas
	if replicas < m.Config.Replication {
		job.Status = Pinning
		job.NextAttempt = m.Now().Add(5 * time.Second)
		return nil
	}
	job.Status = Pinned
	job.NextAttempt = time.Time{}
	job.Error = ""
	if job.PreviousCID != "" && job.PreviousCID != job.CID && !referencedElsewhere(job.PreviousCID, job.Object.ID(), all) {
		if err := m.Pins.Unpin(ctx, job.PreviousCID); err != nil {
			return err
		}
		job.PreviousCID = ""
	}
	return nil
}

func referencedElsewhere(cid, excluded string, jobs []Job) bool {
	for _, job := range jobs {
		if job.Object.ID() == excluded || !job.DesiredPin {
			continue
		}
		// A pending import may produce the same CID. Keep the old pin until
		// its CID is known, even if this delays garbage collection.
		if job.CID == "" {
			return true
		}
		if job.CID == cid || job.PreviousCID == cid {
			return true
		}
	}
	return false
}
