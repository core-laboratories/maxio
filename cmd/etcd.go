// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"context"
	"errors"
	"fmt"

	etcd "go.etcd.io/etcd/client/v3"
)

var errEtcdUnreachable = errors.New("etcd is unreachable, please check your endpoints")

func etcdErrToErr(err error, etcdEndpoints []string) error {
	if err == nil {
		return nil
	}
	switch err {
	case context.DeadlineExceeded:
		return fmt.Errorf("%w %s", errEtcdUnreachable, etcdEndpoints)
	default:
		return fmt.Errorf("unexpected error %w from etcd, please check your endpoints %s", err, etcdEndpoints)
	}
}

func saveKeyEtcdWithTTL(ctx context.Context, client *etcd.Client, key string, data []byte, ttl int64) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()
	lease, err := client.Grant(timeoutCtx, ttl)
	if err != nil {
		return etcdErrToErr(err, client.Endpoints())
	}
	_, err = client.Put(timeoutCtx, key, string(data), etcd.WithLease(lease.ID))
	etcdLogIf(ctx, err)
	return etcdErrToErr(err, client.Endpoints())
}

func saveKeyEtcd(ctx context.Context, client *etcd.Client, key string, data []byte, opts ...options) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()
	if len(opts) > 0 {
		return saveKeyEtcdWithTTL(ctx, client, key, data, opts[0].ttl)
	}
	_, err := client.Put(timeoutCtx, key, string(data))
	etcdLogIf(ctx, err)
	return etcdErrToErr(err, client.Endpoints())
}

func deleteKeyEtcd(ctx context.Context, client *etcd.Client, key string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	_, err := client.Delete(timeoutCtx, key)
	etcdLogIf(ctx, err)
	return etcdErrToErr(err, client.Endpoints())
}

func readKeyEtcd(ctx context.Context, client *etcd.Client, key string) ([]byte, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()
	resp, err := client.Get(timeoutCtx, key)
	if err != nil {
		etcdLogOnceIf(ctx, err, "etcd-retrieve-keys")
		return nil, etcdErrToErr(err, client.Endpoints())
	}
	if resp.Count == 0 {
		return nil, errConfigNotFound
	}
	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			return ev.Value, nil
		}
	}
	return nil, errConfigNotFound
}
