// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package ipfs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"regexp"
	"strings"
)

var cidPattern = regexp.MustCompile(`^b[a-z2-7]{20,}$`)

type Importer interface {
	Add(ctx context.Context, filename string, content io.Reader) (string, error)
	HashOnly(ctx context.Context, filename string, content io.Reader) (string, error)
}

type PinManager interface {
	Pin(ctx context.Context, cid string) error
	Unpin(ctx context.Context, cid string) error
	Replicas(ctx context.Context, cid string) (int, error)
}

type KuboClient struct {
	endpoint string
	token    string
	chunker  string
	client   *http.Client
}

type ClusterClient struct {
	endpoint string
	user     string
	password string
	client   *http.Client
}

func NewKuboClient(c Config) *KuboClient {
	return &KuboClient{c.KuboEndpoint, c.KuboToken, c.Chunker, &http.Client{Timeout: c.RequestTimeout}}
}

func NewClusterClient(c Config) *ClusterClient {
	return &ClusterClient{c.ClusterEndpoint, c.ClusterUser, c.ClusterPassword, &http.Client{Timeout: c.RequestTimeout}}
}

func (k *KuboClient) Add(ctx context.Context, filename string, content io.Reader) (string, error) {
	return k.add(ctx, filename, content, false)
}

func (k *KuboClient) HashOnly(ctx context.Context, filename string, content io.Reader) (string, error) {
	return k.add(ctx, filename, content, true)
}

func (k *KuboClient) add(ctx context.Context, filename string, content io.Reader, onlyHash bool) (string, error) {
	if filename == "" || strings.ContainsAny(filename, "/\\\x00\r\n") {
		return "", errors.New("invalid IPFS import filename")
	}
	u, err := url.Parse(k.endpoint)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, "/api/v0/add")
	q := u.Query()
	q.Set("cid-version", "1")
	q.Set("raw-leaves", "true")
	q.Set("chunker", k.chunker)
	q.Set("pin", "false")
	if onlyHash {
		q.Set("only-hash", "true")
	}
	q.Set("progress", "false")
	u.RawQuery = q.Encode()

	reader, writer := io.Pipe()
	multi := multipart.NewWriter(writer)
	copyDone := make(chan error, 1)
	go func() {
		defer close(copyDone)
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, filename))
		header.Set("Content-Type", "application/octet-stream")
		part, err := multi.CreatePart(header)
		if err == nil {
			_, err = io.Copy(part, content)
		}
		if err == nil {
			err = multi.Close()
		}
		writer.CloseWithError(err)
		copyDone <- err
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), reader)
	if err != nil {
		reader.Close()
		return "", err
	}
	req.Header.Set("Content-Type", multi.FormDataContentType())
	if k.token != "" {
		req.Header.Set("Authorization", "Bearer "+k.token)
	}
	resp, err := k.client.Do(req)
	if err != nil {
		reader.Close()
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		reader.Close()
		return "", fmt.Errorf("Kubo add failed: HTTP %d", resp.StatusCode)
	}
	var result struct {
		Hash string `json:"Hash"`
	}
	dec := json.NewDecoder(io.LimitReader(resp.Body, 1<<20))
	if err := dec.Decode(&result); err != nil {
		reader.Close()
		return "", err
	}
	if err := <-copyDone; err != nil {
		return "", err
	}
	if !cidPattern.MatchString(result.Hash) {
		return "", errors.New("Kubo returned an invalid CIDv1")
	}
	return result.Hash, nil
}

func (c *ClusterClient) request(ctx context.Context, method, target string) (*http.Response, error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return nil, err
	}
	rel, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, rel.Path)
	u.RawQuery = rel.RawQuery
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.password)
	return c.client.Do(req)
}

func (c *ClusterClient) Pin(ctx context.Context, cid string) error {
	if !cidPattern.MatchString(cid) {
		return errors.New("invalid CIDv1")
	}
	resp, err := c.request(ctx, http.MethodPost, "/pins/"+cid+"?replication-min=2&replication-max=2")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("IPFS Cluster pin failed: HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *ClusterClient) Unpin(ctx context.Context, cid string) error {
	if !cidPattern.MatchString(cid) {
		return errors.New("invalid CIDv1")
	}
	resp, err := c.request(ctx, http.MethodDelete, "/pins/"+cid)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("IPFS Cluster unpin failed: HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *ClusterClient) Replicas(ctx context.Context, cid string) (int, error) {
	if !cidPattern.MatchString(cid) {
		return 0, errors.New("invalid CIDv1")
	}
	resp, err := c.request(ctx, http.MethodGet, "/pins/"+cid)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("IPFS Cluster status failed: HTTP %d", resp.StatusCode)
	}
	var status struct {
		PeerMap map[string]struct {
			Status string `json:"status"`
		} `json:"peer_map"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&status); err != nil {
		return 0, err
	}
	count := 0
	for _, peer := range status.PeerMap {
		if strings.EqualFold(peer.Status, "pinned") {
			count++
		}
	}
	return count, nil
}

// ImportVersion identifies the stable UnixFS parameters associated with every CID.
func ImportVersion(c Config) string {
	return "unixfs-cidv1-raw-leaves-sha2-256-" + c.Chunker + "-v1"
}
