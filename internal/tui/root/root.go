package root

import (
	"time"

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

	HeaderModel                *HeaderModel
	BodyModel                  *BodyModel
	loginDialog                *LoginDialog
	CreateNamespaceDialogModel *CreateNamespaceDialogModel
}

func NewRootModel(epinioClient *client.Client, version string) RootModel {
	return RootModel{
		EpinioClient: epinioClient,

		HeaderModel:                NewHeaderModel(version),
		BodyModel:                  NewBodyModel(epinioClient),
		loginDialog:                NewLoginDialog(),
		CreateNamespaceDialogModel: NewCreateNamespaceDialogModel(epinioClient),
	}
}

func doTick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return cmd.TickMsg{}
	})
}

func (m RootModel) Init() tea.Cmd {
	return doTick()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//log.Printf("RootModel/Update, msg: %#v", msg)

	cmds := []tea.Cmd{}

	switch msg := msg.(type) {

	case cmd.TickMsg:
		cmds = append(cmds, doTick())

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

	// update login dialog
	updatedLoginModel, resultMsg := m.loginDialog.Update(msg)
	m.loginDialog = updatedLoginModel.(*LoginDialog)
	cmds = append(cmds, resultMsg)

	// update CreateNamespaceDialogModel
	updatedNamespaceDialogModel, resultMsg := m.CreateNamespaceDialogModel.Update(msg)
	m.CreateNamespaceDialogModel = updatedNamespaceDialogModel.(*CreateNamespaceDialogModel)
	cmds = append(cmds, resultMsg)

	if m.loginDialog.Visible || m.CreateNamespaceDialogModel.Visible {
		return m, tea.Batch(cmds...)
	}

	// update body
	updatedBodyModel, resultMsg := m.BodyModel.Update(msg)
	m.BodyModel = updatedBodyModel.(*BodyModel)
	cmds = append(cmds, resultMsg)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	//log.Printf("RootModel/View")

	if m.loginDialog.Visible {
		dialogView := m.loginDialog.View()
		return lipgloss.Place(10, 20, lipgloss.Center, lipgloss.Center, dialogView)
	}

	if m.CreateNamespaceDialogModel.Visible {
		dialogView := m.CreateNamespaceDialogModel.View()
		return lipgloss.Place(10, 20, lipgloss.Center, lipgloss.Center, dialogView)
	}

	headerView := m.HeaderModel.View()

	body := m.BodyModel.View()

	bodyBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(width - 2).
		Height(height - lipgloss.Height(headerView) - 2).
		Render(body)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		bodyBox,
	)
}
