package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Namespace struct {
	Meta Meta
}

type Meta struct {
	Name      string
	CreatedAt *time.Time
}

type APINamespaceGetMsg struct {
	Result []*Namespace
	Err    error
}

func NewAPINamespaceGetMsg(result []*Namespace, err error) APINamespaceGetMsg {
	return APINamespaceGetMsg{Result: result, Err: err}
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

		var namespaces []*Namespace
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

func NewAPINamespaceDeleteCmd(name string) func() tea.Msg {
	return func() tea.Msg {
		insecureTransport := http.DefaultTransport.(*http.Transport).Clone()
		insecureTransport.TLSClientConfig.InsecureSkipVerify = true

		c := &http.Client{
			Timeout:   10 * time.Second,
			Transport: insecureTransport,
		}

		req, err := http.NewRequest(http.MethodDelete, "https://epinio.172.21.0.4.omg.howdoi.website/api/v1/namespaces/"+name, nil)
		if err != nil {
			return NewAPINamespaceDeleteMsg(err)
		}
		req.SetBasicAuth("admin", "password")

		res, err := c.Do(req)
		if err != nil {
			return NewAPINamespaceDeleteMsg(err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return NewAPINamespaceDeleteMsg(err)
		}
		log.Print(string(b))

		return NewAPINamespaceDeleteMsg(nil)
	}
}
