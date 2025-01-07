// By contributing to, or using this source code, you agree with the terms of the
// MIT-style licensed that can be found below:
//
// Copyright (c) 2025-present Gustavo "Guz" L. de Mello
// Copyright (c) 2025-present The Lored.dev Contributors
// Copyright (c) 2016 The Gitea Authors
// Copyright (c) 2014 The Gogs Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Sections of the contents of this file were sourced from the official Gitea SDK for Go,
// which can be found at https://gitea.com/gitea/go-sdk and is licensed under a MIT-style
// licensed stated at the start of this file.

package gitea

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type client struct {
	endpoint string
	http     *http.Client
}

func newClient(endpoint string, http *http.Client) *client {
	return &client{endpoint, http}
}

func (c *client) GetContents(owner, repo, ref, filepath string) (contentsResponse, error) {
	data, _, err := c.get(
		fmt.Sprintf("/repos/%s/%s/contents/%s?ref=%s", owner, repo, url.QueryEscape(ref), filepath),
	)
	if err != nil {
		return contentsResponse{}, err
	}

	var file contentsResponse
	if err := json.Unmarshal(data, &file); err != nil {
		return contentsResponse{}, errors.Join(
			errors.New("failed to parse JSON response from API"),
			err,
		)
	}

	return file, nil
}

func (c *client) ListContents(owner, repo, ref, filepath string) ([]contentsResponse, error) {
	data, _, err := c.get(
		fmt.Sprintf("/repos/%s/%s/contents/%s?ref=%s", owner, repo, url.QueryEscape(ref), filepath),
	)
	if err != nil {
		return []contentsResponse{}, err
	}

	var directory []contentsResponse
	if err := json.Unmarshal(data, &directory); err != nil {
		return []contentsResponse{}, errors.Join(
			errors.New("failed to parse JSON response from API"),
			err,
		)
	}

	return directory, nil
}

func (c *client) get(path string) (body []byte, res *http.Response, err error) {
	res, err = c.http.Get(c.endpoint + path)
	if err != nil {
		return nil, nil, errors.Join(errors.New("failed to request"), err)
	}
	defer res.Body.Close()

	data, err := statusCodeToErr(res)
	if err != nil {
		return data, res, err
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}

	return data, res, err
}

func statusCodeToErr(resp *http.Response) (body []byte, err error) {
	if resp.StatusCode/100 == 2 {
		return nil, nil
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read on HTTP error %d: %v", resp.StatusCode, err)
	}

	errMap := make(map[string]interface{})
	if err = json.Unmarshal(data, &errMap); err != nil {
		return data, fmt.Errorf(
			"Unknown API Error: %d Request Path: '%s'\nResponse body: '%s'",
			resp.StatusCode,
			resp.Request.URL.Path,
			string(data),
		)
	}

	if msg, ok := errMap["message"]; ok {
		return data, fmt.Errorf("%v", msg)
	}

	return data, fmt.Errorf("%s: %s", resp.Status, string(data))
}

type contentsResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
	SHA  string `json:"sha"`

	// NOTE: can be "file", "dir", "symlink" or "submodule"
	Type string `json:"type"`
	Size int64  `json:"size"`
	// NOTE: populated just when `type` is "contentsResponseTypeFile"
	Encoding *string `json:"encoding"`
	// NOTE: populated just when `type` is "contentsResponseTypeFile"
	Content *string `json:"content"`
	// NOTE: populated just when `type` is "contentsResponseTypeSymlink"
	Target *string `json:"target"`

	URL         *string `json:"url"`
	HTMLURL     *string `json:"html_url"`
	GitURL      *string `json:"git_url"`
	DownloadURL *string `json:"download_url"`

	// NOTE: populated just when `type` is "contentsResponseTypeSubmodule"
	SubmoduleGitURL *string `json:"submodule_giit_url"`

	Links *fileLinksResponse `json:"_links"`
}

type fileLinksResponse struct {
	Self    *string `json:"self"`
	GitURL  *string `json:"git"`
	HTMLURL *string `json:"html"`
}
