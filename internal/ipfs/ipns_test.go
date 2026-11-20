// Copyright (c) 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later

package ipfs

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testIPNSName = "k51qzi5uqu5dlvj2baxnqndepeb86cbk3ng7n3i46uzyxzyqj2xjonzllnv0v8"

func TestSignedIPNSRepublishTag(t *testing.T) {
	policy, err := ParseTags("maxio-ipns-republish=false")
	if err != nil || policy.Republish == nil || *policy.Republish {
		t.Fatalf("republish opt-out: %#v %v", policy, err)
	}
	for _, value := range []string{"true", "always"} {
		if _, err := ParseTags("maxio-ipns-republish=" + value); err == nil {
			t.Fatalf("unsupported republish tag %q accepted", value)
		}
	}
	if _, err := ParseMetadata(map[string]string{"maxio-ipns-republish": "true"}); err == nil {
		t.Fatal("republish metadata accepted")
	}
}

func TestSignedIPNSRepublishStartupDefault(t *testing.T) {
	ref := ObjectRef{Bucket: "shared", Owner: "alice", Key: TenantPrefix("alice") + ".maxio/ipns/" + testIPNSName}
	for _, tc := range []struct {
		startup bool
		optOut  bool
		want    bool
	}{
		{false, false, false},
		{true, false, true},
		{true, true, false},
	} {
		m := &Manager{Config: Config{Enabled: true, IPNSRepublish: tc.startup}, Store: &memoryStore{}, Now: time.Now}
		if err := m.RecordSignedIPNS(context.Background(), ref, tc.optOut); err != nil {
			t.Fatal(err)
		}
		job, _, err := m.Status(context.Background(), ref)
		if err != nil || job.IPNSRepublish != tc.want {
			t.Fatalf("startup=%t optOut=%t: job=%#v err=%v", tc.startup, tc.optOut, job, err)
		}
	}
}

type recordSource struct{ ref ObjectRef }

func (s recordSource) Open(_ context.Context, ref ObjectRef) (io.ReadCloser, error) {
	if ref.ID() == s.ref.ID() {
		return io.NopCloser(strings.NewReader("signed-record")), nil
	}
	return io.NopCloser(strings.NewReader("payload")), nil
}

type fakeSignedPublisher struct {
	inspected, published int
	cid                  string
	expiry               time.Time
}

func (f *fakeSignedPublisher) Inspect(_ context.Context, name string, record []byte) (string, time.Time, error) {
	f.inspected++
	if name != testIPNSName || string(record) != "signed-record" {
		return "", time.Time{}, io.ErrUnexpectedEOF
	}
	if f.expiry.IsZero() {
		f.expiry = time.Now().Add(24 * time.Hour)
	}
	return f.cid, f.expiry, nil
}
func (f *fakeSignedPublisher) Put(_ context.Context, name, cid string, record []byte) error {
	f.published++
	if name != testIPNSName || cid != f.cid || string(record) != "signed-record" {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func TestKuboSignedRecordRequests(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("unexpected request %s", r.URL)
		}
		file, _, err := r.FormFile("record")
		if err != nil {
			t.Error(err)
			return
		}
		data, _ := io.ReadAll(file)
		file.Close()
		if string(data) != "signed-record" {
			t.Errorf("record %q", data)
		}
		switch r.URL.Path {
		case "/api/v0/name/inspect":
			if r.URL.Query().Get("verify") != testIPNSName {
				t.Error("missing name verification")
			}
			json.NewEncoder(w).Encode(map[string]any{"Entry": map[string]string{"Value": "/ipfs/" + testCID, "Validity": time.Now().Add(time.Hour).Format(time.RFC3339Nano)}, "Validation": map[string]any{"Name": testIPNSName, "Valid": true}})
		case "/api/v0/name/put":
			if r.URL.Query().Get("arg") != testIPNSName || r.URL.Query().Get("force") != "false" {
				t.Error("unsafe publication arguments")
			}
			json.NewEncoder(w).Encode(map[string]string{"Name": testIPNSName, "Value": "/ipfs/" + testCID})
		default:
			t.Errorf("unexpected endpoint %s", r.URL.Path)
		}
	}))
	defer server.Close()
	kubo, err := NewKuboClient(Config{KuboEndpoint: server.URL, KuboToken: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	kubo.client = server.Client()
	cid, expiry, err := kubo.Inspect(context.Background(), testIPNSName, []byte("signed-record"))
	if err != nil || cid != testCID || !expiry.After(time.Now()) {
		t.Fatalf("inspect %q %v", cid, err)
	}
	if err := kubo.Put(context.Background(), testIPNSName, cid, []byte("signed-record")); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls %d", calls)
	}
}

func TestKuboCustomCA(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/name/inspect" {
			t.Errorf("unexpected request %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"Entry":      map[string]string{"Value": "/ipfs/" + testCID, "Validity": time.Now().Add(time.Hour).Format(time.RFC3339Nano)},
			"Validation": map[string]any{"Name": testIPNSName, "Valid": true},
		})
	}))
	defer server.Close()
	caFile := filepath.Join(t.TempDir(), "ca.crt")
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	if err := os.WriteFile(caFile, caPEM, 0600); err != nil {
		t.Fatal(err)
	}
	kubo, err := NewKuboClient(Config{KuboEndpoint: server.URL, KuboCAFile: caFile})
	if err != nil {
		t.Fatal(err)
	}
	if cid, _, err := kubo.Inspect(context.Background(), testIPNSName, []byte("record")); err != nil || cid != testCID {
		t.Fatalf("custom Kubo CA: cid=%q err=%v", cid, err)
	}
}

func TestClientSignedIPNSRequiresOwnedPublicPin(t *testing.T) {
	ctx := context.Background()
	owner := "alice"
	target := ObjectRef{Bucket: "shared", Owner: owner, Key: TenantPrefix(owner) + "file", Version: "v1", Size: 7, Modified: time.Unix(1, 0)}
	record := ObjectRef{Bucket: "shared", Owner: owner, Key: TenantPrefix(owner) + ".maxio/ipns/" + testIPNSName, Version: "v1", Modified: time.Unix(2, 0)}
	if SignedRecordName(record) != testIPNSName {
		t.Fatal("record path not recognized")
	}
	store := &memoryStore{}
	names := &fakeSignedPublisher{cid: testCID}
	now := time.Unix(100, 0)
	m := &Manager{Config: Config{Enabled: true, IPNSRepublish: true, Replication: 2, DefaultAccess: Public, GatewayRate: 1 << 20}, Store: store, Source: recordSource{ref: record}, Importer: fakeImporter{}, Pins: &fakePins{replicas: 2}, Names: names, Now: func() time.Time { return now }}
	pin := true
	access := Public
	if err := m.RecordCreated(ctx, target, Policy{Pin: &pin, Access: &access}); err != nil {
		t.Fatal(err)
	}
	if err := m.RecordSignedIPNS(ctx, record, true); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 0 {
		t.Fatal("published before target pin")
	}
	now = now.Add(6 * time.Second)
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 1 {
		t.Fatalf("publication count %d", names.published)
	}
	now = now.Add(ipnsRepublishPeriod)
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 1 {
		t.Fatal("republished without tenant opt-in")
	}
	if err := m.RecordSignedIPNS(ctx, record, false); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 2 {
		t.Fatalf("signed bytes not republished: %d", names.published)
	}
	if err := m.RecordSignedIPNS(ctx, record, true); err != nil {
		t.Fatal(err)
	}
	now = now.Add(ipnsRepublishPeriod)
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 2 {
		t.Fatal("republished after opt-out")
	}
	response := httptest.NewRecorder()
	NewPublicGateway(m).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ipns/"+testIPNSName, nil))
	if response.Code != http.StatusOK || response.Body.String() != "payload" {
		t.Fatalf("gateway %d %q", response.Code, response.Body.String())
	}
	if err := m.RecordSignedIPNS(ctx, record, false); err != nil {
		t.Fatal(err)
	}
	published, _, err := m.Status(ctx, record)
	if err != nil {
		t.Fatal(err)
	}
	published.IPNSExpiry = time.Now().Add(-time.Second)
	names.expiry = published.IPNSExpiry
	if err := store.Save(ctx, published); err != nil {
		t.Fatal(err)
	}
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 2 {
		t.Fatal("expired record was republished")
	}
	response = httptest.NewRecorder()
	NewPublicGateway(m).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ipns/"+testIPNSName, nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("expired name served: %d", response.Code)
	}
	published.IPNSExpiry = time.Now().Add(time.Hour)
	published.Status = Pinned
	if err := store.Save(ctx, published); err != nil {
		t.Fatal(err)
	}
	if err := m.RecordRemoved(ctx, target); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	NewPublicGateway(m).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ipns/"+testIPNSName, nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("removed target served: %d", response.Code)
	}
}

func TestNewSignedRecordVersionSupersedesOld(t *testing.T) {
	ctx := context.Background()
	owner := "alice"
	key := TenantPrefix(owner) + ".maxio/ipns/" + testIPNSName
	old := ObjectRef{Bucket: "shared", Owner: owner, Key: key, Version: "v1", Modified: time.Unix(1, 0)}
	newer := ObjectRef{Bucket: "shared", Owner: owner, Key: key, Version: "v2", Modified: time.Unix(2, 0)}
	store := &memoryStore{}
	target := Job{Object: ObjectRef{Bucket: "shared", Owner: owner, Key: TenantPrefix(owner) + "file", Size: 7}, CID: testCID, DesiredPin: true, Access: Public, Status: Pinned, Replicas: 2, Updated: time.Now()}
	if err := store.Save(ctx, target); err != nil {
		t.Fatal(err)
	}
	names := &fakeSignedPublisher{cid: testCID}
	m := &Manager{Config: Config{Enabled: true, Replication: 2, GatewayRate: 1 << 20}, Store: store, Source: recordSource{ref: newer}, Names: names, Now: time.Now}
	for _, ref := range []ObjectRef{old, newer} {
		if err := m.RecordSignedIPNS(ctx, ref, false); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if names.published != 1 {
		t.Fatalf("expected only latest version to publish, got %d", names.published)
	}
	response := httptest.NewRecorder()
	NewPublicGateway(m).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ipns/"+testIPNSName, nil))
	if response.Code != http.StatusOK {
		t.Fatalf("latest version unavailable: %d", response.Code)
	}
}
