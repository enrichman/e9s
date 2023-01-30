package root

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/constants"
	"github.com/enrichman/e9s/internal/tui/style"
)

type HeaderModel struct {
	version string
}

func NewHeaderModel(version string) HeaderModel {
	return HeaderModel{version: version}
}

func (m HeaderModel) Init() tea.Cmd                           { return nil }
func (m HeaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m HeaderModel) View() string {
	w := lipgloss.Width

	logo := logo()

	infoBoxWrapper := lipgloss.NewStyle().
		Padding(1).
		Width(width - w(logo)).
		Render(infoBox(m.version))

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		infoBoxWrapper,
		logo,
	)
}

func infoBox(version string) string {
	labelBox := lipgloss.JoinVertical(
		lipgloss.Left,
		style.Label.Render("Version:"),
		style.Label.Render("Built with love"),
	)

	valueBox := lipgloss.NewStyle().PaddingLeft(1).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			version,
		),
	)

	return lipgloss.JoinHorizontal(lipgloss.Left, labelBox, valueBox)
}

func logo() string {
	return lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(style.ColorOrange).
		Render(constants.Logo)
}
