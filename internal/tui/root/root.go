package root

import (
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
	BodyModel   tea.Model
	loginDialog tea.Model
}

func NewRootModel(version string) RootModel {
	return RootModel{
		HeaderModel: NewHeaderModel(version),
		BodyModel:   NewBodyModel(),
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

	// handle HttpResults
	case HttpResult:
		{
			if msg.err != nil {
				log.Printf("%+v", msg.err.Error())
			} else {
				// b, err := io.ReadAll(msg.res.Body)
				// log.Printf("%+v %+v", string(b), err)
				// m.HeaderModel.Update(msg)
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

type HttpResult struct {
	res *http.Response
	err error
}

func checkServer() tea.Msg {
	insecureTransport := http.DefaultTransport.(*http.Transport).Clone()
	insecureTransport.TLSClientConfig.InsecureSkipVerify = true

	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: insecureTransport,
	}

	req, err := http.NewRequest(http.MethodGet, "https://epinio.172.21.0.4.omg.howdoi.website/api/v1/namespaces", nil)
	if err != nil {
		return HttpResult{err: err}
	}
	req.SetBasicAuth("admin", "password")

	res, err := c.Do(req)
	if err != nil {
		return HttpResult{err: err}
	}

	return HttpResult{res: res}
}
