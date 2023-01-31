package root

import (
	"io"
	"log"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	width  = 96
	height = 30
)

type RootModel struct {
	HeaderModel tea.Model
	loginDialog tea.Model
}

func NewRootModel(version string) RootModel {
	return RootModel{
		HeaderModel: NewHeaderModel(version),
		loginDialog: NewLoginDialog(),
	}
}

func (m RootModel) Init() tea.Cmd {
	return checkServer
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	log.Printf("%#v", msg)

	// update login dialog
	updatedLoginModel, resultMsg := m.loginDialog.Update(msg)
	m.loginDialog = updatedLoginModel
	cmds = append(cmds, resultMsg)

	switch msg := msg.(type) {

	// update window size
	case tea.WindowSizeMsg:
		{
			width = msg.Width
			height = msg.Height
		}

	// handle HttpResults
	case HttpResult:
		{
			if msg.err != nil {
				log.Print("error!!")
			} else {
				b, err := io.ReadAll(msg.res.Body)
				log.Printf("%+v %+v", string(b), err)
				m.HeaderModel.Update(msg)
			}
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

	headerView := m.HeaderModel.View()

	body := m.loginDialog.View()

	if body == "" {
		body = m.body()
	}

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

const url = "https://charm.sh/"

type HttpResult struct {
	res *http.Response
	err error
}

func checkServer() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(url)
	if err != nil {
		return HttpResult{err: err}
	}
	defer res.Body.Close()

	return HttpResult{res: res}
}
