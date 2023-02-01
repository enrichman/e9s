package root

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BodyModel struct {
	Namespaces []*Namespace

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

type Namespace struct {
	Meta Meta
}

type Meta struct {
	Name      string
	CreatedAt *time.Time
}

func (m *BodyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// handle HttpResults
	case HttpResult:
		{
			if msg.err != nil {
				log.Printf("%+v", msg.err.Error())
			} else {
				b, err := io.ReadAll(msg.res.Body)
				if err != nil {
					log.Printf("%+v", err.Error())
					break
				}
				log.Print(string(b))

				var namespaces []*Namespace
				err = json.Unmarshal(b, &namespaces)
				if err != nil {
					log.Printf("%+v", err.Error())
					break
				}

				m.Namespaces = namespaces

				columns := []table.Column{
					{Title: "NAME", Width: 20},
					{Title: "CREATED", Width: 40},
					{Title: "APPLICATIONS", Width: 14},
					{Title: "CONFIGURATIONS", Width: width - 84},
				}
				m.table.SetColumns(columns)

				rows := []table.Row{}
				for _, ns := range namespaces {
					rows = append(rows, table.Row{ns.Meta.Name, ns.Meta.CreatedAt.String(), "", ""})
				}

				m.table.SetHeight(6)

				m.table.SetRows(rows)
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			m.table.MoveUp(1)
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			m.table.MoveDown(1)

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selectedRow := m.table.SelectedRow()
			log.Print(selectedRow[0])

		}
	}

	return m, nil
}

func (m *BodyModel) View() string {
	if len(m.Namespaces) > 0 {
		return m.table.View() + "\n"
	}
	return "Press q to quit.\n"
}
