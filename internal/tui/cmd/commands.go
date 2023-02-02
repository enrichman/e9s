package cmd

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enrichman/e9s/pkg/client"
)

type TickMsg struct{}
type ShowCreateNamespaceDialogMsg struct{}

type APINamespaceGetStartMsg struct{}
type APINamespaceDeleteStartMsg struct{ Name string }

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

type APINamespaceGetResultMsg struct {
	Result []*client.Namespace
	Err    error
}

func NewAPINamespaceGetCmd(namespaceService NamespaceService) func() tea.Msg {
	return func() tea.Msg {
		log.Printf("Executing NamespaceGetCmd")
		namespaces, err := namespaceService.Get(context.Background())
		if err != nil {
			log.Printf("Error during NamespaceGetCmd: %v", err.Error())
		}
		return APINamespaceGetResultMsg{Result: namespaces, Err: err}
	}
}

type APINamespaceDeleteResultMsg struct {
	DeletedNamespace string
	Err              error
}

func NewAPINamespaceDeleteCmd(namespaceService NamespaceService, name string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Executing NamespaceDeleteCmd [%s]", name)
		_, err := namespaceService.Delete(context.Background(), name)
		if err != nil {
			log.Printf("Error during NamespaceDeleteCmd: %v", err.Error())
		}
		return APINamespaceDeleteResultMsg{DeletedNamespace: name, Err: err}
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
