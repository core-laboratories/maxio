// Copyright (c) 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later

package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

const maxSignedRecordSize = 10 << 10

var ipnsNamePattern = regexp.MustCompile(`^k[0-9a-z]{20,127}$`)

// SignedRecordPublisher validates and broadcasts a client-signed IPNS record.
// Neither MaxIO nor Kubo needs the signing key.
type SignedRecordPublisher interface {
	Inspect(ctx context.Context, name string, record []byte) (string, time.Time, error)
	Put(ctx context.Context, name, cid string, record []byte) error
}

// SignedRecordName identifies the tenant-owned S3 control object for an IPNS name.
func SignedRecordName(ref ObjectRef) string {
	if !ref.Valid() {
		return ""
	}
	prefix := TenantPrefix(ref.Owner) + ".maxio/ipns/"
	if !strings.HasPrefix(ref.Key, prefix) {
		return ""
	}
	name := strings.TrimPrefix(ref.Key, prefix)
	if !ipnsNamePattern.MatchString(name) {
		return ""
	}
	return name
}

func (k *KuboClient) ipnsRecordRequest(ctx context.Context, command, name string, record []byte, result any) error {
	if !ipnsNamePattern.MatchString(name) || len(record) == 0 || len(record) > maxSignedRecordSize {
		return errors.New("invalid signed IPNS record")
	}
	u, err := url.Parse(k.endpoint)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "api/v0/name/"+command)
	q := u.Query()
	if command == "inspect" {
		q.Set("verify", name)
		q.Set("dump", "false")
	} else {
		q.Set("arg", name)
		q.Set("force", "false")
	}
	u.RawQuery = q.Encode()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("record", "record.ipns")
	if err != nil {
		return err
	}
	if _, err := part.Write(record); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if k.token != "" {
		req.Header.Set("Authorization", "Bearer "+k.token)
	}
	resp, err := k.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Kubo name/%s failed: HTTP %d", command, resp.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(result)
}

func (k *KuboClient) Inspect(ctx context.Context, name string, record []byte) (string, time.Time, error) {
	var inspected struct {
		Entry struct {
			Value    string `json:"Value"`
			Validity string `json:"Validity"`
		} `json:"Entry"`
		Validation struct {
			Name  string `json:"Name"`
			Valid bool   `json:"Valid"`
		} `json:"Validation"`
	}
	if err := k.ipnsRecordRequest(ctx, "inspect", name, record, &inspected); err != nil {
		return "", time.Time{}, err
	}
	cid := strings.TrimPrefix(inspected.Entry.Value, "/ipfs/")
	if !inspected.Validation.Valid || inspected.Validation.Name != name || cid == inspected.Entry.Value || !cidPattern.MatchString(cid) {
		return "", time.Time{}, errors.New("invalid IPNS signature or target CID")
	}
	expiry, err := time.Parse(time.RFC3339Nano, inspected.Entry.Validity)
	if err != nil || !expiry.After(time.Now()) {
		return "", time.Time{}, errors.New("expired or invalid IPNS validity")
	}
	return cid, expiry, nil
}

func (k *KuboClient) Put(ctx context.Context, name, cid string, record []byte) error {
	if !cidPattern.MatchString(cid) {
		return errors.New("invalid IPNS target CID")
	}
	var published struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	}
	if err := k.ipnsRecordRequest(ctx, "put", name, record, &published); err != nil {
		return err
	}
	if published.Name != name || published.Value != "/ipfs/"+cid {
		return errors.New("Kubo returned an unexpected IPNS publication")
	}
	return nil
}

func (m *Manager) processSignedIPNS(ctx context.Context, job *Job, all []Job) error {
	if !job.DesiredPin {
		job.Status = Unpinned
		job.IPNSName = ""
		job.IPNSCID = ""
		job.CID = ""
		return nil
	}
	name := SignedRecordName(job.Object)
	if name == "" || m.Names == nil {
		return errors.New("signed IPNS publisher is unavailable")
	}
	reader, err := m.Source.Open(ctx, job.Object)
	if err != nil {
		return err
	}
	record, readErr := io.ReadAll(io.LimitReader(reader, maxSignedRecordSize+1))
	closeErr := reader.Close()
	if readErr != nil {
		return readErr
	}
	if closeErr != nil {
		return closeErr
	}
	if len(record) == 0 || len(record) > maxSignedRecordSize {
		return errors.New("signed IPNS record exceeds size limit")
	}
	cid, expiry, err := m.Names.Inspect(ctx, name, record)
	if err != nil {
		return err
	}
	if !expiry.After(time.Now()) {
		return errors.New("signed IPNS record has expired")
	}
	job.CID = cid
	for _, target := range all {
		if target.SignedIPNS || target.Object.Owner != job.Object.Owner || target.Object.Bucket != job.Object.Bucket || target.CID != cid || !target.DesiredPin || target.Access != Public || target.Status != Pinned || target.Replicas < m.Config.Replication {
			continue
		}
		if err := m.Names.Put(ctx, name, cid, record); err != nil {
			return err
		}
		job.IPNSName = name
		job.IPNSCID = cid
		job.IPNSUpdated = m.Now()
		job.IPNSExpiry = expiry
		job.Replicas = target.Replicas
		job.Status = Pinned
		job.NextAttempt = time.Time{}
		return nil
	}
	job.Status = Pinning
	job.NextAttempt = m.Now().Add(5 * time.Second)
	return nil
}
