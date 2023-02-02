package root

import (
	"log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/tui/cmd"
	"github.com/enrichman/e9s/pkg/client"
)

type Body struct {
	EpinioClient *client.Client
	Namespaces   []*client.Namespace

	table      table.Model
	loadingGet bool
}

func NewBody(epinioClient *client.Client) Body {

	t := table.New(
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(false)
	//s.Cell = s.Cell.Foreground(lipgloss.Color("12"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("12"))

	t.SetStyles(s)

	return Body{
		EpinioClient: epinioClient,
		table:        t,
	}
}

func updateBody(body Body, msg tea.Msg) (Body, tea.Cmd) {
	//log.Printf("BodyModel/Update, msg: %#v", msg)

	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case cmd.TickMsg:
		{
			// check if is already loading
			if !body.loadingGet {
				cmds = append(cmds, cmd.NewAPINamespaceGetCmd(body.EpinioClient.Namespaces))
				body.loadingGet = true
			}
		}

	case cmd.APINamespaceGetMsg:
		{
			body.loadingGet = false // GET ended
			body.Namespaces = msg.Result
			body.table = updateTable(body.table, body.Namespaces)
		}

	case cmd.APINamespaceDeleteMsg:
		{
			if msg.Err != nil {
				log.Printf("%+v", msg.Err.Error())
				break
			}

			// check if is already loading
			if !body.loadingGet {
				cmds = append(cmds, cmd.NewAPINamespaceGetCmd(body.EpinioClient.Namespaces))
				body.loadingGet = true
			}
		}

	case cmd.APINamespaceCreateMsg:
		{
			if msg.Err != nil {
				log.Printf("%+v", msg.Err.Error())
				break
			}

			// check if is already loading
			if !body.loadingGet {
				cmds = append(cmds, cmd.NewAPINamespaceGetCmd(body.EpinioClient.Namespaces))
				body.loadingGet = true
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			body.table.MoveUp(1)
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			body.table.MoveDown(1)

		case "ctrl+d":
			selectedNamespace := body.table.SelectedRow()[0]
			cmds = append(cmds, cmd.NewAPINamespaceDeleteCmd(body.EpinioClient.Namespaces, selectedNamespace))
		}
	}

	return body, tea.Batch(cmds...)
}

func updateTable(namespaceTable table.Model, namespaces []*client.Namespace) table.Model {
	columns := []table.Column{
		{Title: "NAME", Width: 20},
		{Title: "CREATED", Width: 40},
		{Title: "APPLICATIONS", Width: 14},
		{Title: "CONFIGURATIONS", Width: width - 84},
	}
	namespaceTable.SetColumns(columns)

	rows := []table.Row{}
	for _, ns := range namespaces {
		rows = append(rows, table.Row{ns.Meta.Name, ns.Meta.CreatedAt.String(), "", ""})
	}
	namespaceTable.SetRows(rows)

	namespaceTable.SetHeight(6)

	return namespaceTable
}

func viewBody(body Body) string {
	if len(body.Namespaces) == 0 {
		return ""
	}
	return viewTable(body.table)
}

func viewTable(table table.Model) string {
	return table.View()
}
