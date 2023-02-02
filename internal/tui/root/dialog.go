package root

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginDialog struct {
	Visible              bool
	input                textinput.Model
	cancelButtonSelected bool
	okButtonSelected     bool
}

func NewLoginDialog() *LoginDialog {
	ti := textinput.New()
	ti.Placeholder = "http://epinio.example.com"
	ti.Focus()
	ti.CharLimit = 200

	return &LoginDialog{input: ti}
}

func (m *LoginDialog) Init() tea.Cmd {
	return nil
}

func (m *LoginDialog) Update(msg tea.Msg) tea.Cmd {
	//log.Printf("loginDialog/Update, msg: %#v", msg)

	cmds := []tea.Cmd{}

	if !m.Visible {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+l":
				m.Visible = true
			}
		}
		return nil
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
			m.Visible = false

		case "left", "right":
			m.toggle()

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
		}
	}
	return tea.Batch(cmds...)
}

func (d *LoginDialog) toggle() {
	d.cancelButtonSelected = !d.cancelButtonSelected
	d.okButtonSelected = !d.okButtonSelected
}

func (d *LoginDialog) View() string {
	if !d.Visible {
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
