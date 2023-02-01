package cmd

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enrichman/e9s/pkg/client"
)

type NamespaceService interface {
	Get(ctx context.Context) ([]*client.Namespace, error)
	Delete(ctx context.Context, name string) (client.OkResponse, error)
}

type APINamespaceGetMsg struct {
	Result []*client.Namespace
	Err    error
}

func NewAPINamespaceGetMsg(result []*client.Namespace, err error) APINamespaceGetMsg {
	return APINamespaceGetMsg{Result: result, Err: err}
}

func NewAPINamespaceGetCmd(namespaceService NamespaceService) func() tea.Msg {
	return func() tea.Msg {
		log.Printf("Executing NamespaceGetCmd")
		namespaces, err := namespaceService.Get(context.Background())
		return NewAPINamespaceGetMsg(namespaces, err)
	}
}

type APINamespaceDeleteMsg struct {
	Err error
}

func NewAPINamespaceDeleteMsg(err error) APINamespaceDeleteMsg {
	return APINamespaceDeleteMsg{Err: err}
}

func NewAPINamespaceDeleteCmd(namespaceService NamespaceService, name string) func() tea.Msg {
	return func() tea.Msg {
		log.Printf("Executing NamespaceDeleteCmd [%s]", name)
		_, err := namespaceService.Delete(context.Background(), name)
		return NewAPINamespaceDeleteMsg(err)
	}
}
