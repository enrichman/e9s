package client

import (
	"context"
	"time"
)

type NamespaceService struct {
	client *Client
}

type Namespace struct {
	Meta Meta
}

type Meta struct {
	Name      string
	CreatedAt *time.Time
}

type OkResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}

type APIError struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

func (s *NamespaceService) Get(ctx context.Context) ([]*Namespace, error) {
	namespaces, err := get(ctx, s.client, "/namespaces", []*Namespace{})
	return namespaces, err
}

func (s *NamespaceService) Delete(ctx context.Context, name string) (OkResponse, error) {
	res, err := delete(ctx, s.client, "/namespaces/"+name, OkResponse{})
	return res, err
}

func (s *NamespaceService) Create(ctx context.Context, name string) (OkResponse, error) {
	type namespaceCreateRequest struct {
		Name string `json:"name"`
	}

	payload := &namespaceCreateRequest{
		Name: name,
	}

	okResponse, err := post(ctx, s.client, "/namespaces", payload, OkResponse{})
	return okResponse, err
}
