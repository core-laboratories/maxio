// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/core-laboratories/maxio/internal/auth"
	"github.com/core-laboratories/maxio/internal/event"
	"github.com/core-laboratories/maxio/internal/ipfs"
	"github.com/minio/pkg/v3/policy"
)

const ipfsJobPrefix = "ipfs/jobs/"
const ipfsTenantPrefix = "ipfs/tenants/"

var globalIPFSManager *ipfs.Manager
var globalIPFSMu sync.RWMutex
var ipfsKnownTenants sync.Map

type ipfsTenantRecord struct {
	ID        string    `json:"id"`
	Principal string    `json:"principal"`
	Created   time.Time `json:"created"`
}

type ipfsObjectStore struct{ api ObjectLayer }

func (s ipfsObjectStore) Load(ctx context.Context, id string) (ipfs.Job, bool, error) {
	data, err := readConfig(ctx, s.api, ipfsJobPrefix+id+".json")
	if errors.Is(err, errConfigNotFound) {
		return ipfs.Job{}, false, nil
	}
	if err != nil {
		return ipfs.Job{}, false, err
	}
	var job ipfs.Job
	if err := json.Unmarshal(data, &job); err != nil {
		return ipfs.Job{}, false, err
	}
	return job, true, nil
}

func (s ipfsObjectStore) Save(ctx context.Context, job ipfs.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return saveConfig(ctx, s.api, ipfsJobPrefix+job.Object.ID()+".json", data)
}

func (s ipfsObjectStore) CompareAndSwap(ctx context.Context, expected uint64, job ipfs.Job) (bool, error) {
	id := job.Object.ID()
	lock := s.api.NewNSLock(minioMetaBucket, ipfsJobPrefix+id+".json")
	lkctx, err := lock.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return false, err
	}
	defer lock.Unlock(lkctx)
	current, found, err := s.Load(lkctx.Context(), id)
	if err != nil {
		return false, err
	}
	if (found && current.Revision != expected) || (!found && expected != 0) {
		return false, nil
	}
	job.Revision = expected + 1
	return true, s.Save(lkctx.Context(), job)
}

func (s ipfsObjectStore) List(ctx context.Context) ([]ipfs.Job, error) {
	items := make(chan itemOrErr[ObjectInfo], 100)
	if err := s.api.Walk(ctx, minioMetaBucket, ipfsJobPrefix, items, WalkOptions{}); err != nil {
		return nil, err
	}
	var jobs []ipfs.Job
	for item := range items {
		if item.Err != nil {
			return nil, item.Err
		}
		id := strings.TrimSuffix(strings.TrimPrefix(item.Item.Name, ipfsJobPrefix), ".json")
		job, ok, err := s.Load(ctx, id)
		if err != nil {
			return nil, err
		}
		if ok {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func (s ipfsObjectStore) Open(ctx context.Context, ref ipfs.ObjectRef) (io.ReadCloser, error) {
	reader, err := s.api.GetObjectNInfo(ctx, ref.Bucket, ref.Key, nil, http.Header{}, ObjectOptions{VersionID: ref.Version})
	if err != nil {
		return nil, err
	}
	if ref.ETag != "" && reader.ObjInfo.ETag != ref.ETag {
		reader.Close()
		return nil, errors.New("IPFS object changed since job creation")
	}
	if !ref.Modified.IsZero() && !reader.ObjInfo.ModTime.Equal(ref.Modified) {
		reader.Close()
		return nil, errors.New("IPFS object changed since job creation")
	}
	return reader, nil
}

func initIPFS(api ObjectLayer) error {
	c, err := ipfs.ConfigFromEnv()
	if err != nil {
		return err
	}
	if !c.Enabled {
		return nil
	}
	globalIPFSMu.Lock()
	defer globalIPFSMu.Unlock()
	if globalIPFSManager != nil {
		return nil
	}
	store := ipfsObjectStore{api: api}
	manager := ipfs.NewManager(c, store, store)
	globalIPFSManager = manager
	if c.PublicGateway {
		server := &http.Server{Addr: c.GatewayListen, Handler: ipfs.NewPublicGateway(manager), ReadHeaderTimeout: 10 * time.Second}
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				loggerIPFSError(err)
			}
		}()
		go func() {
			<-GlobalContext.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		}()
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		reconcile := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		defer reconcile.Stop()
		if err := ipfsWithWorkerLock(GlobalContext, store, func(ctx context.Context) error {
			return ipfsReconcile(ctx, manager, store)
		}); err != nil {
			loggerIPFSError(err)
		}
		for {
			select {
			case <-GlobalContext.Done():
				return
			case <-ticker.C:
				if err := ipfsWithWorkerLock(GlobalContext, store, manager.RunOnce); err != nil {
					loggerIPFSError(err)
				}
			case <-reconcile.C:
				if err := ipfsWithWorkerLock(GlobalContext, store, func(ctx context.Context) error {
					return ipfsReconcile(ctx, manager, store)
				}); err != nil {
					loggerIPFSError(err)
				}
			}
		}
	}()
	return nil
}

func ipfsWithWorkerLock(ctx context.Context, store ipfsObjectStore, work func(context.Context) error) error {
	lock := store.api.NewNSLock(minioMetaBucket, "ipfs/worker")
	lkctx, err := lock.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	defer lock.Unlock(lkctx)
	return work(lkctx.Context())
}

func ipfsReconcile(ctx context.Context, manager *ipfs.Manager, store ipfsObjectStore) error {
	owners := make(map[string]string)
	tenantItems := make(chan itemOrErr[ObjectInfo], 100)
	if err := store.api.Walk(ctx, minioMetaBucket, ipfsTenantPrefix, tenantItems, WalkOptions{LatestOnly: true}); err != nil {
		return err
	}
	for item := range tenantItems {
		if item.Err != nil {
			return item.Err
		}
		data, err := readConfig(ctx, store.api, item.Item.Name)
		if err != nil {
			return err
		}
		var record ipfsTenantRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		if record.ID == ipfs.TenantID(record.Principal) {
			owners[record.ID] = record.Principal
		}
	}
	objects := make(chan itemOrErr[ObjectInfo], 100)
	if err := store.api.Walk(ctx, manager.Config.TenantBucket, "users/", objects, WalkOptions{}); err != nil {
		return err
	}
	seen := make(map[string]bool)
	for item := range objects {
		if item.Err != nil {
			return item.Err
		}
		info := item.Item
		parts := strings.SplitN(info.Name, "/", 3)
		if len(parts) != 3 || parts[0] != "users" || info.DeleteMarker {
			continue
		}
		owner := owners[parts[1]]
		if owner == "" || !ipfs.AuthorizedKey(owner, info.Name) {
			continue
		}
		ref := ipfs.ObjectRef{Bucket: manager.Config.TenantBucket, Key: info.Name, Version: info.VersionID, Owner: owner, ETag: info.ETag, Size: info.Size, Modified: info.ModTime}
		seen[ref.ID()] = true
		job, found, err := manager.Status(ctx, ref)
		if err != nil {
			return err
		}
		if !found || job.Object.ETag != ref.ETag || !job.Object.Modified.Equal(ref.Modified) {
			// A missed event has no trustworthy authorization marker for copied controls.
			// Recover the object record without publishing its bytes.
			pin := false
			if err := manager.RecordCreated(ctx, ref, ipfs.Policy{Pin: &pin}); err != nil {
				return err
			}
		}
	}
	jobs, err := store.List(ctx)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		if job.Object.Bucket == manager.Config.TenantBucket && !seen[job.Object.ID()] && job.DesiredPin {
			if err := manager.RecordRemoved(ctx, job.Object); err != nil {
				return err
			}
		}
	}
	return nil
}

func loggerIPFSError(err error) {
	// The error must not include credentials or object metadata.
	fmt.Printf("IPFS background operation: %v\n", err)
}

func ipfsEvent(args eventArgs) {
	globalIPFSMu.RLock()
	manager := globalIPFSManager
	globalIPFSMu.RUnlock()
	if manager == nil || args.Object.Name == "" || args.BucketName == "" {
		return
	}
	if args.BucketName != manager.Config.TenantBucket {
		return
	}
	if args.ReqParams["principalId"] != "" && !ipfs.AuthorizedKey(args.ReqParams["principalId"], args.Object.Name) {
		return
	}
	ref := ipfs.ObjectRef{
		Bucket: args.BucketName, Key: args.Object.Name,
		Version: args.Object.VersionID, Owner: args.ReqParams["principalId"],
		ETag: args.Object.ETag, Size: args.Object.Size, Modified: args.Object.ModTime,
	}
	ctx, cancel := context.WithTimeout(GlobalContext, 15*time.Second)
	defer cancel()
	switch args.EventName {
	case event.ObjectCreatedPut, event.ObjectCreatedCopy, event.ObjectCreatedCompleteMultipartUpload,
		event.ObjectCreatedPutTagging, event.ObjectCreatedDeleteTagging, event.ObjectRestoreCompleted:
		var base, tags ipfs.Policy
		if args.IPFSMetaControl {
			var err error
			base, err = ipfs.ParseMetadata(args.Object.UserDefined)
			if err != nil {
				loggerIPFSError(err)
				return
			}
		}
		if args.IPFSTagControl {
			var err error
			tags, err = ipfs.ParseTags(args.Object.UserTags)
			if err != nil {
				loggerIPFSError(err)
				return
			}
		}
		replaceBase := args.EventName != event.ObjectCreatedPutTagging && args.EventName != event.ObjectCreatedDeleteTagging
		if err := manager.RecordPolicy(ctx, ref, base, tags, replaceBase); err != nil {
			loggerIPFSError(err)
		}
	case event.ObjectRemovedDelete, event.ObjectRemovedDeleteAllVersions:
		if err := manager.RecordRemoved(ctx, ref); err != nil {
			loggerIPFSError(err)
		}
	}
}

func validateIPFSControls(ctx context.Context, r *http.Request, bucket, object string, metadata map[string]string, tagString string) APIErrorCode {
	globalIPFSMu.RLock()
	manager := globalIPFSManager
	globalIPFSMu.RUnlock()
	if manager == nil || bucket != manager.Config.TenantBucket {
		return ErrNone
	}
	meta, err := ipfs.ParseMetadata(metadata)
	if err != nil {
		return ErrInvalidRequest
	}
	tags, err := ipfs.ParseTags(tagString)
	if err != nil {
		return ErrInvalidRequest
	}
	if meta.IPNS != "" || meta.IPNSAuto != nil || len(meta.CoreIDs) > 0 || tags.IPNS != "" || tags.IPNSAuto != nil {
		return ErrInvalidRequest
	}
	if tags.Pin != nil {
		meta.Pin = tags.Pin
	}
	if tags.Access != nil {
		meta.Access = tags.Access
	}
	access := manager.Config.DefaultAccess
	if meta.Access != nil {
		access = *meta.Access
	}
	if ipfs.EffectivePin(manager.Config.AutoPin, meta) && access != ipfs.Public {
		return ErrInvalidRequest
	}
	return ErrNone
}

func ipfsTenantAllowed(ctx context.Context, r *http.Request, action policy.Action, bucket, object string, cred auth.Credentials, owner bool) bool {
	globalIPFSMu.RLock()
	manager := globalIPFSManager
	globalIPFSMu.RUnlock()
	if manager == nil || bucket != manager.Config.TenantBucket || owner {
		return true
	}
	principal := cred.AccessKey
	if cred.ParentUser != "" {
		principal = cred.ParentUser
	}
	if principal == "" {
		return false
	}
	if err := ipfsEnsureTenant(ctx, manager, principal); err != nil {
		loggerIPFSError(err)
		return false
	}
	if object != "" {
		return ipfs.AuthorizedKey(principal, object)
	}
	if r.Method == http.MethodHead {
		return true
	}
	if action == policy.ListBucketAction || action == policy.ListBucketVersionsAction {
		return strings.HasPrefix(r.URL.Query().Get("prefix"), ipfs.TenantPrefix(principal))
	}
	return true
}

func ipfsEnsureTenant(ctx context.Context, manager *ipfs.Manager, principal string) error {
	id := ipfs.TenantID(principal)
	if _, ok := ipfsKnownTenants.Load(id); ok {
		return nil
	}
	store, ok := manager.Store.(ipfsObjectStore)
	if !ok {
		return errors.New("IPFS tenant store is unavailable")
	}
	name := ipfsTenantPrefix + id + ".json"
	lock := store.api.NewNSLock(minioMetaBucket, name)
	lkctx, err := lock.GetLock(ctx, globalOperationTimeout)
	if err != nil {
		return err
	}
	defer lock.Unlock(lkctx)
	data, err := readConfig(lkctx.Context(), store.api, name)
	if err == nil {
		var record ipfsTenantRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return err
		}
		if record.ID != id || record.Principal != principal {
			return errors.New("tenant identity conflict")
		}
		ipfsKnownTenants.Store(id, principal)
		return nil
	}
	if !errors.Is(err, errConfigNotFound) {
		return err
	}
	record := ipfsTenantRecord{ID: id, Principal: principal, Created: time.Now().UTC()}
	data, err = json.Marshal(record)
	if err != nil {
		return err
	}
	if err := saveConfig(lkctx.Context(), store.api, name, data); err != nil {
		return err
	}
	ipfsKnownTenants.Store(id, principal)
	return nil
}

func ipfsStatusHeaders(ctx context.Context, r *http.Request, w http.ResponseWriter, info ObjectInfo) {
	globalIPFSMu.RLock()
	manager := globalIPFSManager
	globalIPFSMu.RUnlock()
	if manager == nil || !manager.Config.StatusHeaders || info.Bucket != manager.Config.TenantBucket || getRequestAuthType(r) == authTypeAnonymous {
		return
	}
	cred := getReqAccessCred(r, globalSite.Region())
	principal := cred.AccessKey
	if cred.ParentUser != "" {
		principal = cred.ParentUser
	}
	if !ipfs.AuthorizedKey(principal, info.Name) {
		return
	}
	ref := ipfs.ObjectRef{Bucket: info.Bucket, Key: info.Name, Version: info.VersionID, Owner: principal}
	job, found, err := manager.Status(ctx, ref)
	if err != nil {
		loggerIPFSError(err)
		return
	}
	if !found || job.Object.ETag != info.ETag || !job.Object.Modified.Equal(info.ModTime) {
		return
	}
	w.Header().Set("x-maxio-pin-status", string(job.Status))
	w.Header().Set("x-maxio-replicas", fmt.Sprintf("%d", job.Replicas))
	if job.CID != "" {
		w.Header().Set("x-maxio-cid", job.CID)
	}
}

func filterIPFSPrivateMetadata(info *ObjectInfo) {
	globalIPFSMu.RLock()
	enabled := globalIPFSManager != nil
	globalIPFSMu.RUnlock()
	if !enabled {
		return
	}
	info.UserDefined = cloneMSS(info.UserDefined)
	for key := range info.UserDefined {
		if strings.EqualFold(key, "x-amz-meta-maxio-core-ids") || strings.EqualFold(key, "maxio-core-ids") {
			delete(info.UserDefined, key)
		}
	}
}
