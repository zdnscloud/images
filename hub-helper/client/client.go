package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	baseURL     = "https://hub.docker.com"
	registryURL = "https://registry-1.docker.io"
	loginPath   = "/v2/users/login/"

	defaultPageSize = 100
)

type Client struct {
	client     *http.Client
	namespace  string
	token      string
	credential LoginCredential
}

func NewClient(user, password string) (*Client, error) {
	client := &Client{
		namespace: user,
		client: &http.Client{
			Transport: http.DefaultTransport,
		},
	}

	client.credential = LoginCredential{
		User:     user,
		Password: password,
	}
	err := client.refreshToken()
	if err != nil {
		return nil, fmt.Errorf("login to dockerhub error: %v", err)
	}

	return client, nil
}

func (c *Client) refreshToken() error {
	b, err := json.Marshal(c.credential)
	if err != nil {
		return fmt.Errorf("marshal credential error: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+loginPath, bytes.NewReader(b))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("login to dockerhub error: %s", string(body))
	}

	token := &TokenResp{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return fmt.Errorf("unmarshal token response error: %v", err)
	}

	c.token = token.Token
	return nil
}

func (c *Client) DeleteTag(repo, tag string) error {
	resp, err := c.Do(http.MethodDelete, deleteTagPath(c.namespace, repo, tag), nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("%d -- %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) GetRepos(page, pageSize int) (*ReposResp, error) {
	resp, err := c.Do(http.MethodGet, listReposPath(c.namespace, page, pageSize), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%d -- %s", resp.StatusCode, string(body))
	}

	repos := &ReposResp{}
	err = json.Unmarshal(body, repos)
	if err != nil {
		return nil, fmt.Errorf("unmarshal repos list %s error: %v", string(body), err)
	}

	return repos, nil
}

func (c *Client) GetTags(repo string, page, pageSize int) (*TagsResp, error) {
	resp, err := c.Do(http.MethodGet, listTagsPath(c.namespace, repo, page, pageSize), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("%d -- %s", resp.StatusCode, body)
	}

	tags := &TagsResp{}
	err = json.Unmarshal(body, tags)
	if err != nil {
		return nil, fmt.Errorf("unmarshal tags list %s error: %v", string(body), err)
	}

	return tags, nil
}

func (c *Client) Do(method, path string, body io.Reader) (*http.Response, error) {
	url := baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil || method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", c.token))

	return c.client.Do(req)
}

func (c *Client) Namespce() string {
	return c.namespace
}

func listReposPath(namespace string, page, pageSize int) string {
	return fmt.Sprintf("/v2/repositories/%s/?page=%d&page_size=%d", namespace, page, pageSize)
}

func listTagsPath(namespace, repo string, page, pageSize int) string {
	return fmt.Sprintf("/v2/repositories/%s/%s/tags/?page=%d&page_size=%d", namespace, repo, page, pageSize)
}

func deleteTagPath(namespace, repo, tag string) string {
	return fmt.Sprintf("/v2/repositories/%s/%s/tags/%s/", namespace, repo, tag)
}
