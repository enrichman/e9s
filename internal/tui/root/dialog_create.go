package root

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/tui/cmd"
	"github.com/enrichman/e9s/pkg/client"
)

type CreateNamespaceDialog struct {
	EpinioClient         *client.Client
	Visible              bool
	input                textinput.Model
	cancelButtonSelected bool
	okButtonSelected     bool
}

func NewCreateNamespaceDialog(epinioClient *client.Client) *CreateNamespaceDialog {
	ti := textinput.New()

	ti.CharLimit = 50

	return &CreateNamespaceDialog{
		EpinioClient: epinioClient,
		input:        ti,
	}
}

func (m *CreateNamespaceDialog) Update(msg tea.Msg) tea.Cmd {
	//log.Printf("CreateNamespaceDialogModel/Update, msg: %+v", msg)
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {

	// key press
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			m.Visible = false
			m.input.Reset()

		case "left", "right":
			m.toggle()

		case "enter":
			cmds = append(cmds, cmd.NewAPINamespaceCreateCmd(m.EpinioClient.Namespaces, m.input.Value()))
			m.Visible = false
			m.input.Reset()
		}

	case cmd.ShowCreateNamespaceDialogMsg:
		m.Visible = true
		m.input.Reset()
		m.input.Focus()
		return nil
	}

	updatedInput, inputCmd := m.input.Update(msg)
	m.input = updatedInput
	cmds = append(cmds, inputCmd)

	return tea.Batch(cmds...)
}

func (d *CreateNamespaceDialog) toggle() {
	d.cancelButtonSelected = !d.cancelButtonSelected
	d.okButtonSelected = !d.okButtonSelected
}

func (d *CreateNamespaceDialog) View() string {
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

	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Create namespace")
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
