package cmd

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enrichman/e9s/pkg/client"
)

type TickMsg struct{}
type ShowCreateNamespaceDialogMsg struct{}

func NewCmd(msg tea.Msg) func() tea.Msg {
	return func() tea.Msg {
		return msg
	}
}

type NamespaceService interface {
	Get(ctx context.Context) ([]*client.Namespace, error)
	Delete(ctx context.Context, name string) (client.OkResponse, error)
	Create(ctx context.Context, name string) (client.OkResponse, error)
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
		if err != nil {
			log.Printf("Error during NamespaceGetCmd: %v", err.Error())
		}
		return NewAPINamespaceGetMsg(namespaces, err)
	}
}

type APINamespaceDeleteMsg struct {
	Err error
}

func NewAPINamespaceDeleteMsg(err error) APINamespaceDeleteMsg {
	return APINamespaceDeleteMsg{Err: err}
}

func NewAPINamespaceDeleteCmd(namespaceService NamespaceService, name string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Executing NamespaceDeleteCmd [%s]", name)
		_, err := namespaceService.Delete(context.Background(), name)
		if err != nil {
			log.Printf("Error during NamespaceDeleteCmd: %v", err.Error())
		}
		return NewAPINamespaceDeleteMsg(err)
	}
}

type APINamespaceCreateMsg struct {
	Err error
}

func NewAPINamespaceCreateMsg(err error) APINamespaceCreateMsg {
	return APINamespaceCreateMsg{Err: err}
}

func NewAPINamespaceCreateCmd(namespaceService NamespaceService, name string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Executing NamespaceCreateCmd")
		_, err := namespaceService.Create(context.Background(), name)
		if err != nil {
			log.Printf("Error during NamespaceCreateCmd: %v", err.Error())
		}
		return NewAPINamespaceCreateMsg(err)
	}
}
