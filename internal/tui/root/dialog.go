package root

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loginDialog struct {
	visible              bool
	input                textinput.Model
	cancelButtonSelected bool
	okButtonSelected     bool
}

func NewLoginDialog() *loginDialog {
	ti := textinput.New()
	ti.Placeholder = "http://epinio.example.com"
	ti.Focus()
	ti.CharLimit = 200

	return &loginDialog{input: ti}
}

func (m *loginDialog) Init() tea.Cmd {
	return nil
}

func (m *loginDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("loginDialog/Update, msg: %#v", msg)

	cmds := []tea.Cmd{}

	if !m.visible {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+l":
				m.visible = true
			}
		}
		return m, nil
	}

	updatedInput, cmd := m.input.Update(msg)
	m.input = updatedInput
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	// key press
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
		case "down", "j":

		case "esc":
			m.visible = false

		case "left", "right":
			m.toggle()

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			log.Print("enter")
		}
	}
	return m, tea.Batch(cmds...)
}

func (d *loginDialog) toggle() {
	d.cancelButtonSelected = !d.cancelButtonSelected
	d.okButtonSelected = !d.okButtonSelected
}

func (d *loginDialog) View() string {
	if !d.visible {
		return ""
	}

	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		Margin(2, 1, 0)

	activeButtonStyle := buttonStyle.Copy().
		Background(lipgloss.Color("#F25D94")).
		Underline(true)

	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1)

	okButton := ""
	cancelButton := ""
	if d.okButtonSelected {
		okButton = activeButtonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("Maybe")
	} else {
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("Maybe")
	}

	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)

	ui := lipgloss.JoinVertical(
		lipgloss.Center,
		question,
		lipgloss.NewStyle().Width(50).Align(lipgloss.Left).Render(d.input.View()),
		buttons,
	)

	dialog := lipgloss.Place(
		width-2, 0,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	return dialog
}
