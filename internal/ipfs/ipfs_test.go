// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const testCID = "bafkreigh2akiscaildcboe65umrn2t7w4j3s2w2tnwqg5eao5c43y6h3ay"

func TestPolicies(t *testing.T) {
	meta, err := ParseMetadata(map[string]string{"X-Amz-Meta-Maxio-Pin": "true", "X-Amz-Meta-Maxio-Access": "public"})
	if err != nil || meta.Pin == nil || !*meta.Pin || *meta.Access != Public {
		t.Fatalf("metadata: %#v %v", meta, err)
	}
	tags, err := ParseTags("maxio-pin=false&ordinary=value")
	if err != nil || tags.Pin == nil || *tags.Pin {
		t.Fatalf("tags: %#v %v", tags, err)
	}
	if EffectivePin(false, meta, tags) {
		t.Fatal("tag must override metadata")
	}
	for _, bad := range []map[string]string{
		{"maxio-pin": "yes"},
		{"maxio-access": "internal"},
		{"maxio-ipns": "../escape"},
		{"maxio-unknown": "value"},
	} {
		if _, err := ParseMetadata(bad); err == nil {
			t.Fatalf("accepted invalid metadata: %#v", bad)
		}
	}
}

func TestKuboAddStreamsCIDv1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/add" || r.URL.Query().Get("cid-version") != "1" || r.URL.Query().Get("pin") != "false" || r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected request %s", r.URL)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Error(err)
			return
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil || string(data) != "hello" {
			t.Errorf("content %q %v", data, err)
		}
		json.NewEncoder(w).Encode(map[string]string{"Hash": testCID})
	}))
	defer server.Close()
	c := &KuboClient{endpoint: server.URL, token: "test-token", chunker: "size-262144", client: server.Client()}
	cid, err := c.Add(context.Background(), "hello.txt", bytes.NewBufferString("hello"))
	if err != nil || cid != testCID {
		t.Fatalf("add: %s %v", cid, err)
	}
}

func TestClusterReplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != "user" || p != "pass" {
			t.Error("missing cluster authentication")
		}
		if r.Method == http.MethodPost {
			if r.URL.Query().Get("replication-min") != "2" || r.URL.Query().Get("replication-max") != "2" {
				t.Error("wrong replication factor")
			}
			return
		}
		if r.Method == http.MethodGet {
			io.WriteString(w, `{"peer_map":{"peer1":{"status":"pinned"},"peer2":{"status":"pinned"}}}`)
		}
	}))
	defer server.Close()
	c := &ClusterClient{endpoint: server.URL, user: "user", password: "pass", client: server.Client()}
	if err := c.Pin(context.Background(), testCID); err != nil {
		t.Fatal(err)
	}
	count, err := c.Replicas(context.Background(), testCID)
	if err != nil || count != 2 {
		t.Fatalf("replicas: %d %v", count, err)
	}
}

type memoryStore struct {
	sync.Mutex
	jobs map[string]Job
}

func (s *memoryStore) Load(_ context.Context, id string) (Job, bool, error) {
	s.Lock()
	defer s.Unlock()
	j, ok := s.jobs[id]
	return j, ok, nil
}
func (s *memoryStore) Save(_ context.Context, j Job) error {
	s.Lock()
	defer s.Unlock()
	if s.jobs == nil {
		s.jobs = make(map[string]Job)
	}
	s.jobs[j.Object.ID()] = j
	return nil
}
func (s *memoryStore) CompareAndSwap(_ context.Context, expected uint64, j Job) (bool, error) {
	s.Lock()
	defer s.Unlock()
	if s.jobs == nil {
		s.jobs = make(map[string]Job)
	}
	current, found := s.jobs[j.Object.ID()]
	if (found && current.Revision != expected) || (!found && expected != 0) {
		return false, nil
	}
	j.Revision = expected + 1
	s.jobs[j.Object.ID()] = j
	return true, nil
}
func (s *memoryStore) List(_ context.Context) ([]Job, error) {
	s.Lock()
	defer s.Unlock()
	var out []Job
	for _, j := range s.jobs {
		out = append(out, j)
	}
	return out, nil
}

type fakeSource struct{}

func (fakeSource) Open(context.Context, ObjectRef) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("payload")), nil
}

type fakeImporter struct{}

func (fakeImporter) Add(context.Context, string, io.Reader) (string, error)      { return testCID, nil }
func (fakeImporter) HashOnly(context.Context, string, io.Reader) (string, error) { return testCID, nil }

type fakePins struct{ pinned, unpinned int }

func (f *fakePins) Pin(context.Context, string) error             { f.pinned++; return nil }
func (f *fakePins) Unpin(context.Context, string) error           { f.unpinned++; return nil }
func (f *fakePins) Replicas(context.Context, string) (int, error) { return 2, nil }

func TestReferenceCountedUnpin(t *testing.T) {
	store := &memoryStore{}
	pins := &fakePins{}
	m := &Manager{Config: Config{Enabled: true, AutoPin: false, Replication: 2, DefaultAccess: Public, Chunker: "size-262144"}, Store: store, Source: fakeSource{}, Importer: fakeImporter{}, Pins: pins}
	m.Now = func() time.Time { return time.Unix(1, 0) }
	yes := true
	a := Public
	policy := Policy{Pin: &yes, Access: &a}
	first := ObjectRef{Bucket: "shared", Key: "users/a/one", Version: "v1", Owner: "a", ETag: "e1"}
	second := ObjectRef{Bucket: "shared", Key: "users/b/two", Version: "v2", Owner: "b", ETag: "e2"}
	for _, ref := range []ObjectRef{first, second} {
		if err := m.RecordCreated(context.Background(), ref, policy); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, ref := range []ObjectRef{first, second} {
		job, ok, _ := m.Status(context.Background(), ref)
		if !ok || job.Status != Pinned || job.Replicas != 2 {
			t.Fatalf("job %#v", job)
		}
	}
	if err := m.RecordRemoved(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if pins.unpinned != 0 {
		t.Fatal("unpin removed shared content")
	}
	if err := m.RecordRemoved(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if pins.unpinned != 1 {
		t.Fatalf("expected one global unpin, got %d", pins.unpinned)
	}
}

func TestPublicGatewaySeparatesPrivateContent(t *testing.T) {
	store := &memoryStore{}
	ref := ObjectRef{Bucket: "shared", Key: "users/a/file", Version: "v1", Owner: "a", Size: 7}
	job := Job{Object: ref, DesiredPin: true, Access: Private, CID: testCID, Status: Pinned, Replicas: 2}
	if err := store.Save(context.Background(), job); err != nil {
		t.Fatal(err)
	}
	m := &Manager{Config: Config{Enabled: true, Replication: 2, GatewayRate: 1 << 20}, Store: store, Source: fakeSource{}}
	gateway := NewPublicGateway(m)
	request := httptest.NewRequest(http.MethodGet, "/ipfs/"+testCID, nil)
	record := httptest.NewRecorder()
	gateway.ServeHTTP(record, request)
	if record.Code != http.StatusNotFound {
		t.Fatalf("private content returned %d", record.Code)
	}
	job.Access = Public
	if err := store.Save(context.Background(), job); err != nil {
		t.Fatal(err)
	}
	request = httptest.NewRequest(http.MethodGet, "/ipfs/"+testCID, nil)
	request.Header.Set("Range", "bytes=1-3")
	record = httptest.NewRecorder()
	gateway.ServeHTTP(record, request)
	if record.Code != http.StatusPartialContent || record.Body.String() != "ayl" {
		t.Fatalf("range: %d %q", record.Code, record.Body.String())
	}
	if record.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing safe content header")
	}
}

func TestTenantPrefixIsolation(t *testing.T) {
	alice := TenantPrefix("alice")
	bob := TenantPrefix("bob")
	if alice == bob || !AuthorizedKey("alice", alice+"file") {
		t.Fatal("tenant prefix mapping failed")
	}
	for _, key := range []string{bob + "file", "users/client-selected/file", alice[:len(alice)-1] + "-other/file"} {
		if AuthorizedKey("alice", key) {
			t.Fatalf("cross-tenant key accepted: %q", key)
		}
	}
}

func TestConfigRejectsUnsafePrivatePin(t *testing.T) {
	c := Config{Enabled: true, AutoPin: true, DefaultAccess: Private, Replication: 2, TenantBucket: "shared", Chunker: "size-262144"}
	if err := c.Validate(); err == nil {
		t.Fatal("private plaintext auto pin accepted")
	}
	c.DefaultAccess = Public
	c.ClusterEndpoint = "https://cluster.example"
	c.KuboEndpoint = "https://kubo.example"
	c.KuboToken = "secret"
	c.ClusterUser = "user"
	c.ClusterPassword = "secret"
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	c.KuboEndpoint = "http://kubo.example"
	if err := c.Validate(); err == nil {
		t.Fatal("insecure remote Kubo RPC accepted")
	}
}

func TestPublicGatewayRequiresVerifiedReplicas(t *testing.T) {
	store := &memoryStore{}
	ref := ObjectRef{Bucket: "shared", Key: "users/a/file", Version: "v1", Owner: "a", Size: 7}
	job := Job{Object: ref, DesiredPin: true, Access: Public, CID: testCID, Status: Pinning, Replicas: 1}
	if err := store.Save(context.Background(), job); err != nil {
		t.Fatal(err)
	}
	m := &Manager{Config: Config{Enabled: true, Replication: 2, GatewayRate: 1 << 20}, Store: store, Source: fakeSource{}}
	gateway := NewPublicGateway(m)
	request := httptest.NewRequest(http.MethodGet, "/ipfs/"+testCID, nil)
	for _, status := range []JobStatus{Pinning, Failed, Pinned} {
		job.Status = status
		if err := store.Save(context.Background(), job); err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()
		gateway.ServeHTTP(response, request)
		if response.Code != http.StatusNotFound {
			t.Fatalf("unverified status %s returned %d", status, response.Code)
		}
	}
}

func TestTagPinOverrideAndRecovery(t *testing.T) {
	store := &memoryStore{}
	pins := &fakePins{}
	m := &Manager{Config: Config{Enabled: true, AutoPin: false, Replication: 2, DefaultAccess: Public, Chunker: "size-262144"}, Store: store, Source: fakeSource{}, Importer: fakeImporter{}, Pins: pins, Now: time.Now}
	ref := ObjectRef{Bucket: "shared", Key: "users/a/file", Version: "v1", Owner: "a", ETag: "e1"}
	if err := m.RecordCreated(context.Background(), ref, Policy{}); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	job, _, _ := m.Status(context.Background(), ref)
	if job.Status != Unpinned {
		t.Fatalf("default status: %s", job.Status)
	}
	yes := true
	if err := m.RecordPolicy(context.Background(), ref, Policy{}, Policy{Pin: &yes}, false); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	job, _, _ = m.Status(context.Background(), ref)
	if job.Status != Pinned || pins.pinned != 1 {
		t.Fatalf("tag pin: %#v, pin calls %d", job, pins.pinned)
	}
	no := false
	if err := m.RecordPolicy(context.Background(), ref, Policy{}, Policy{Pin: &no}, false); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	job, _, _ = m.Status(context.Background(), ref)
	if job.Status != Unpinned || pins.unpinned != 1 {
		t.Fatalf("tag unpin: %#v, unpin calls %d", job, pins.unpinned)
	}
}

func TestOverwriteWithSameETagRequeues(t *testing.T) {
	store := &memoryStore{}
	m := &Manager{Config: Config{Enabled: true, Replication: 2, DefaultAccess: Public}, Store: store, Now: time.Now}
	ref := ObjectRef{Bucket: "shared", Key: "users/a/file", Owner: "a", ETag: "same", Modified: time.Unix(10, 0)}
	if err := m.RecordCreated(context.Background(), ref, Policy{}); err != nil {
		t.Fatal(err)
	}
	first, _, _ := m.Status(context.Background(), ref)
	ref.Modified = time.Unix(20, 0)
	if err := m.RecordCreated(context.Background(), ref, Policy{}); err != nil {
		t.Fatal(err)
	}
	second, _, _ := m.Status(context.Background(), ref)
	if second.Revision <= first.Revision || !second.Object.Modified.Equal(ref.Modified) {
		t.Fatalf("overwrite not requeued: before %#v after %#v", first, second)
	}
}

type selectivelyMissingSource struct{}

func (selectivelyMissingSource) Open(_ context.Context, ref ObjectRef) (io.ReadCloser, error) {
	if strings.Contains(ref.Key, "missing") {
		return nil, io.EOF
	}
	return io.NopCloser(strings.NewReader("payload")), nil
}

func TestPublicGatewayFallsBackToValidReference(t *testing.T) {
	store := &memoryStore{}
	for _, key := range []string{"users/a/missing", "users/b/present"} {
		ref := ObjectRef{Bucket: "shared", Key: key, Owner: key, Size: 7}
		if err := store.Save(context.Background(), Job{Object: ref, DesiredPin: true, Access: Public, CID: testCID, Status: Pinned, Replicas: 2}); err != nil {
			t.Fatal(err)
		}
	}
	m := &Manager{Config: Config{Enabled: true, Replication: 2, GatewayRate: 1 << 20}, Store: store, Source: selectivelyMissingSource{}}
	response := httptest.NewRecorder()
	NewPublicGateway(m).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ipfs/"+testCID, nil))
	if response.Code != http.StatusOK || response.Body.String() != "payload" {
		t.Fatalf("valid duplicate reference not served: %d %q", response.Code, response.Body.String())
	}
}

func TestIPFSEnabledByDefault(t *testing.T) {
	t.Setenv("MAXIO_IPFS_ENABLED", "")
	t.Setenv("MAXIO_IPFS_AUTO_PIN", "")
	t.Setenv("MAXIO_IPFS_TENANT_BUCKET", "shared")
	t.Setenv("MAXIO_IPFS_CLUSTER_ENDPOINT", "https://cluster.example")
	t.Setenv("MAXIO_IPFS_CLUSTER_USER", "user")
	t.Setenv("MAXIO_IPFS_CLUSTER_PASSWORD", "secret")
	t.Setenv("MAXIO_IPFS_KUBO_ENDPOINT", "https://kubo.example")
	t.Setenv("MAXIO_IPFS_KUBO_TOKEN", "secret")
	config, err := ConfigFromEnv()
	if err != nil || !config.Enabled || config.AutoPin {
		t.Fatalf("default configuration: enabled=%v autoPin=%v err=%v", config.Enabled, config.AutoPin, err)
	}
	t.Setenv("MAXIO_IPFS_ENABLED", "false")
	config, err = ConfigFromEnv()
	if err != nil || config.Enabled {
		t.Fatalf("explicit disable: enabled=%v err=%v", config.Enabled, err)
	}
}
