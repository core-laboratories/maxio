// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package subnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	xhttp "github.com/core-laboratories/maxio/internal/http"
)

const (
	respBodyLimit = 1 << 20 // 1 MiB

	// LoggerWebhookName - subnet logger webhook target
	LoggerWebhookName = "subnet"
)

// Upload given file content (payload) to specified URL
func (c Config) Upload(reqURL string, filename string, payload []byte) (string, error) {
	if !c.Registered() {
		return "", errors.New("Deployment is not registered with SUBNET. Please register the deployment via 'mc license register ALIAS'")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, e := writer.CreateFormFile("file", filename)
	if e != nil {
		return "", e
	}

	if _, e = part.Write(payload); e != nil {
		return "", e
	}
	writer.Close()

	r, e := http.NewRequest(http.MethodPost, reqURL, &body)
	if e != nil {
		return "", e
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())

	return c.submitPost(r)
}

func (c Config) submitPost(r *http.Request) (string, error) {
	configLock.RLock()
	r.Header.Set(xhttp.SubnetAPIKey, c.APIKey)
	configLock.RUnlock()
	r.Header.Set(xhttp.MinioDeploymentID, xhttp.GlobalDeploymentID)

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: c.transport,
	}

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer xhttp.DrainBody(resp.Body)

	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, respBodyLimit))
	if err != nil {
		return "", err
	}
	respStr := string(respBytes)

	if resp.StatusCode == http.StatusOK {
		return respStr, nil
	}

	return respStr, fmt.Errorf("SUBNET request failed with code %d and error: %s", resp.StatusCode, respStr)
}

// Post submit 'payload' to specified URL
func (c Config) Post(reqURL string, payload any) (string, error) {
	if !c.Registered() {
		return "", errors.New("Deployment is not registered with SUBNET. Please register the deployment via 'mc license register ALIAS'")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	r, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	r.Header.Set("Content-Type", "application/json")

	return c.submitPost(r)
}
