package root

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/tui/cmd"
	"github.com/enrichman/e9s/pkg/client"
)

var (
	width  = 96
	height = 30
)

type RootModel struct {
	EpinioClient *client.Client

	HeaderModel tea.Model
	BodyModel   tea.Model
	loginDialog tea.Model
}

func NewRootModel(epinioClient *client.Client, version string) RootModel {
	return RootModel{
		EpinioClient: epinioClient,

		HeaderModel: NewHeaderModel(version),
		BodyModel:   NewBodyModel(epinioClient),
		loginDialog: NewLoginDialog(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return cmd.NewAPINamespaceGetCmd(m.EpinioClient.Namespaces)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("RootModel/Update, msg: %#v", msg)

	cmds := []tea.Cmd{}

	// update login dialog
	updatedLoginModel, resultMsg := m.loginDialog.Update(msg)
	m.loginDialog = updatedLoginModel
	cmds = append(cmds, resultMsg)

	// update body
	updatedBodyModel, resultMsg := m.BodyModel.Update(msg)
	m.BodyModel = updatedBodyModel
	cmds = append(cmds, resultMsg)

	switch msg := msg.(type) {

	// update window size
	case tea.WindowSizeMsg:
		{
			width = msg.Width
			height = msg.Height
		}

	// key press
	case tea.KeyMsg:
		switch msg.String() {

		// exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":

		// The "down" and "j" keys move the cursor down
		case "down", "j":

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":

		}
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	dialogView := m.loginDialog.View()
	if dialogView != "" {
		return dialogView
	}

	headerView := m.HeaderModel.View()

	body := m.BodyModel.View()

	bodyBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(width - 2).
		Height(height - lipgloss.Height(headerView) - 2).
		Render(body)

	lipgloss.Place(10, 10, lipgloss.Center, lipgloss.Center, dialogView)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		bodyBox,
	)
}
