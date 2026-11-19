// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	xioutil "github.com/core-laboratories/maxio/internal/ioutil"
)

const (
	retryInterval = 3 * time.Second
)

type logger = func(ctx context.Context, err error, id string, errKind ...any)

// ErrNotConnected - indicates that the target connection is not active.
var ErrNotConnected = errors.New("not connected to target server/service")

// Target - store target interface
type Target interface {
	Name() string
	SendFromStore(key Key) error
}

// Store - Used to persist items.
type Store[I any] interface {
	Put(item I) (Key, error)
	PutMultiple(item []I) (Key, error)
	Get(key Key) (I, error)
	GetMultiple(key Key) ([]I, error)
	GetRaw(key Key) ([]byte, error)
	PutRaw(b []byte) (Key, error)
	Len() int
	List() []Key
	Del(key Key) error
	Open() error
	Delete() error
}

// Key denotes the key present in the store.
type Key struct {
	Name      string
	Compress  bool
	Extension string
	ItemCount int
}

// String returns the filepath name
func (k Key) String() string {
	keyStr := k.Name
	if k.ItemCount > 1 {
		keyStr = fmt.Sprintf("%d:%s", k.ItemCount, k.Name)
	}
	return keyStr + k.Extension + func() string {
		if k.Compress {
			return compressExt
		}
		return ""
	}()
}

func getItemCount(k string) (count int, err error) {
	count = 1
	v := strings.Split(k, ":")
	if len(v) == 2 {
		return strconv.Atoi(v[0])
	}
	return count, err
}

func parseKey(k string) (key Key) {
	key.Name = k
	if strings.HasSuffix(k, compressExt) {
		key.Compress = true
		key.Name = strings.TrimSuffix(key.Name, compressExt)
	}
	if key.ItemCount, _ = getItemCount(k); key.ItemCount > 1 {
		key.Name = strings.TrimPrefix(key.Name, fmt.Sprintf("%d:", key.ItemCount))
	}
	if vals := strings.Split(key.Name, "."); len(vals) == 2 {
		key.Extension = "." + vals[1]
		key.Name = strings.TrimSuffix(key.Name, key.Extension)
	}
	return key
}

// replayItems - Reads the items from the store and replays.
func replayItems[I any](store Store[I], doneCh <-chan struct{}, log logger, id string) <-chan Key {
	keyCh := make(chan Key)

	go func() {
		defer xioutil.SafeClose(keyCh)

		retryTicker := time.NewTicker(retryInterval)
		defer retryTicker.Stop()

		for {
			for _, key := range store.List() {
				select {
				case keyCh <- key:
				// Get next key.
				case <-doneCh:
					return
				}
			}

			select {
			case <-retryTicker.C:
			case <-doneCh:
				return
			}
		}
	}()

	return keyCh
}

// sendItems - Reads items from the store and re-plays.
func sendItems(target Target, keyCh <-chan Key, doneCh <-chan struct{}, logger logger) {
	retryTicker := time.NewTicker(retryInterval)
	defer retryTicker.Stop()

	send := func(key Key) bool {
		for {
			err := target.SendFromStore(key)
			if err == nil {
				break
			}

			logger(
				context.Background(),
				fmt.Errorf("unable to send log entry to '%s' err '%w'", target.Name(), err),
				target.Name(),
			)

			select {
			// Retrying after 3secs back-off
			case <-retryTicker.C:
			case <-doneCh:
				return false
			}
		}
		return true
	}

	for {
		select {
		case key, ok := <-keyCh:
			if !ok {
				return
			}

			if !send(key) {
				return
			}
		case <-doneCh:
			return
		}
	}
}

// StreamItems reads the keys from the store and replays the corresponding item to the target.
func StreamItems[I any](store Store[I], target Target, doneCh <-chan struct{}, logger logger) {
	go func() {
		keyCh := replayItems(store, doneCh, logger, target.Name())
		sendItems(target, keyCh, doneCh, logger)
	}()
}
