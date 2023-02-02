package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	BaseURL    string
	Username   string
	Password   string

	Namespaces *NamespaceService
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	c := &Client{
		httpClient: httpClient,
		BaseURL:    baseURL,
	}

	c.Namespaces = &NamespaceService{c}

	return c
}

func get[R any](ctx context.Context, c *Client, path string, response R) (R, error) {
	return do(ctx, c, http.MethodGet, path, nil, response)
}

func post[R any](ctx context.Context, c *Client, path string, payload any, response R) (R, error) {
	return do(ctx, c, http.MethodPost, path, payload, response)
}

func delete[R any](ctx context.Context, c *Client, path string, response R) (R, error) {
	return do(ctx, c, http.MethodDelete, path, nil, response)
}

func do[R any](ctx context.Context, c *Client, method, path string, payload any, response R) (R, error) {
	emptyData := new(R)

	body := new(bytes.Buffer)
	if payload != nil {
		err := json.NewEncoder(body).Encode(payload)
		if err != nil {
			return *emptyData, err
		}
	}

	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return *emptyData, err
	}

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return *emptyData, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return *emptyData, err
	}
	log.Print(string(b))

	if res.StatusCode > 399 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(b, &errorResponse)
		if err != nil {
			return *emptyData, err
		}
		return *emptyData, errors.New(errorResponse.Errors[0].Title)
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return *emptyData, err
	}

	return response, nil
}
