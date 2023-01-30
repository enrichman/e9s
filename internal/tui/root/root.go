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
	choices     []string         // items on the to-do list
	cursor      int              // which to-do list item our cursor is pointing at
	selected    map[int]struct{} // which to-do items are selected
}

func NewRootModel(version string) RootModel {
	return RootModel{
		HeaderModel: NewHeaderModel(version),
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m RootModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return checkServer
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("%#v", msg)

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
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated RootModel to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m RootModel) View() string {

	headerView := m.HeaderModel.View()

	bodyBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(width - 2).
		Height(height - lipgloss.Height(headerView) - 2).
		Render(m.body())

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
