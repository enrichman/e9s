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

	Header                Header
	Body                  Body
	LoginDialog           *LoginDialog
	CreateNamespaceDialog *CreateNamespaceDialog
}

func NewRootModel(epinioClient *client.Client, version string) RootModel {
	return RootModel{
		EpinioClient: epinioClient,

		Header:                NewHeader(version),
		Body:                  NewBody(epinioClient),
		LoginDialog:           NewLoginDialog(),
		CreateNamespaceDialog: NewCreateNamespaceDialog(epinioClient),
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
		{
			cmds = append(cmds, doTick())
		}

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

		case "c":
			cmds = append(cmds, cmd.NewCmd(cmd.ShowCreateNamespaceDialogMsg{}))
		}
	}

	// update login dialog
	cmds = append(cmds, m.LoginDialog.Update(msg))

	// update CreateNamespaceDialogModel
	cmds = append(cmds, m.CreateNamespaceDialog.Update(msg))

	if m.LoginDialog.Visible || m.CreateNamespaceDialog.Visible {
		return m, tea.Batch(cmds...)
	}

	// update body
	updatedBody, resultCmds := updateBody(m.Body, msg)
	m.Body = updatedBody
	cmds = append(cmds, resultCmds)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	//log.Printf("RootModel/View")

	if m.LoginDialog.Visible {
		dialogView := m.LoginDialog.View()
		return lipgloss.Place(10, 20, lipgloss.Center, lipgloss.Center, dialogView)
	}

	if m.CreateNamespaceDialog.Visible {
		dialogView := m.CreateNamespaceDialog.View()
		return lipgloss.Place(10, 20, lipgloss.Center, lipgloss.Center, dialogView)
	}

	headerView := viewHeader(m.Header)

	body := viewBody(m.Body)

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
