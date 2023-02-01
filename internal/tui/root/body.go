package root

import (
	"log"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/tui/cmd"
	"github.com/enrichman/e9s/pkg/client"
)

type BodyModel struct {
	Namespaces []*client.Namespace

	table table.Model
}

func NewBodyModel() *BodyModel {

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

	return &BodyModel{table: t}
}

func (m *BodyModel) Init() tea.Cmd {
	return nil
}

func (m *BodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("BodyModel/Update, msg: %#v", msg)

	cmds := []tea.Cmd{}

	switch msg := msg.(type) {

	case cmd.APINamespaceGetMsg:
		{
			m.Namespaces = msg.Result

			columns := []table.Column{
				{Title: "NAME", Width: 20},
				{Title: "CREATED", Width: 40},
				{Title: "APPLICATIONS", Width: 14},
				{Title: "CONFIGURATIONS", Width: width - 84},
			}
			m.table.SetColumns(columns)

			rows := []table.Row{}
			for _, ns := range m.Namespaces {
				rows = append(rows, table.Row{ns.Meta.Name, ns.Meta.CreatedAt.String(), "", ""})
			}
			m.table.SetRows(rows)

			m.table.SetHeight(6)
		}

	case cmd.APINamespaceDeleteMsg:
		{
			if msg.Err != nil {
				log.Printf("%+v", msg.Err.Error())
				break
			}
			cmds = append(cmds, cmd.NewAPINamespaceGetCmd())
		}

	case tea.KeyMsg:
		switch msg.String() {
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			m.table.MoveUp(1)
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			m.table.MoveDown(1)

		case tea.KeyCtrlD.String():
			insecureTransport := http.DefaultTransport.(*http.Transport).Clone()
			insecureTransport.TLSClientConfig.InsecureSkipVerify = true

			c := &http.Client{
				Timeout:   10 * time.Second,
				Transport: insecureTransport,
			}

			ep := client.NewClient(c, "https://epinio.172.21.0.4.omg.howdoi.website/api/v1")
			ep.Username = "admin"
			ep.Password = "password"

			selectedNamespace := m.table.SelectedRow()[0]
			cmds = append(cmds, cmd.NewAPINamespaceDeleteCmd(ep.Namespace, selectedNamespace))

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selectedRow := m.table.SelectedRow()
			log.Print(selectedRow[0])

		}
	}

	return m, tea.Batch(cmds...)
}

func (m *BodyModel) View() string {
	if len(m.Namespaces) > 0 {
		return m.table.View() + "\n"
	}
	return "Press q to quit.\n"
}
