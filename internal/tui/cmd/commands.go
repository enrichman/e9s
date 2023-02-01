package cmd

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enrichman/e9s/pkg/client"
)

type APINamespaceGetMsg struct {
	Result []*client.Namespace
	Err    error
}

func NewAPINamespaceGetMsg(result []*client.Namespace, err error) APINamespaceGetMsg {
	return APINamespaceGetMsg{Result: result, Err: err}
}

type NamespaceDeleter interface {
	Delete(ctx context.Context, name string) (client.OkResponse, error)
}

func NewAPINamespaceGetCmd() func() tea.Msg {
	return func() tea.Msg {
		insecureTransport := http.DefaultTransport.(*http.Transport).Clone()
		insecureTransport.TLSClientConfig.InsecureSkipVerify = true

		c := &http.Client{
			Timeout:   10 * time.Second,
			Transport: insecureTransport,
		}

		req, err := http.NewRequest(http.MethodGet, "https://epinio.172.21.0.4.omg.howdoi.website/api/v1/namespaces", nil)
		if err != nil {
			return NewAPINamespaceGetMsg(nil, err)
		}
		req.SetBasicAuth("admin", "password")

		res, err := c.Do(req)
		if err != nil {
			return NewAPINamespaceGetMsg(nil, err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return NewAPINamespaceGetMsg(nil, err)
		}

		var namespaces []*client.Namespace
		err = json.Unmarshal(b, &namespaces)
		if err != nil {
			return NewAPINamespaceGetMsg(nil, err)
		}

		return NewAPINamespaceGetMsg(namespaces, nil)
	}
}

type APINamespaceDeleteMsg struct {
	Err error
}

func NewAPINamespaceDeleteMsg(err error) APINamespaceDeleteMsg {
	return APINamespaceDeleteMsg{Err: err}
}

func NewAPINamespaceDeleteCmd(deleter NamespaceDeleter, name string) func() tea.Msg {
	return func() tea.Msg {
		log.Printf("Executing NamespaceDeleteCmd [%s]", name)
		_, err := deleter.Delete(context.Background(), name)
		return NewAPINamespaceDeleteMsg(err)
	}
}
