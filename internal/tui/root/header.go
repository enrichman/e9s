package root

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HeaderModel struct {
	version string
}

func NewHeaderModel(version string) HeaderModel {
	return HeaderModel{version: version}
}

func (m HeaderModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m HeaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Println(msg)
	return m, nil
}

const AppName = `___________________       
\_   _____/   __   \______
 |    __)_\____    /  ___/
 |        \  /    /\___ \ 
/_______  / /____//____  >
        \/             \/ 
`

func (m HeaderModel) View() string {
	w := lipgloss.Width

	firstBoxStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	leftBox := lipgloss.JoinVertical(
		lipgloss.Left,
		firstBoxStyle.Render("Version:"),
		firstBoxStyle.Render("Version2sdadsd:"),
	)

	leftBox2 := lipgloss.NewStyle().PaddingLeft(1).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.version,
			m.version+"jnsdkncks",
		),
	)

	rightLogo := lipgloss.NewStyle().Padding(0, 2).Render(AppName)

	boxWrapper := lipgloss.NewStyle().
		Padding(1).
		Width(width - w(rightLogo)).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left, leftBox, leftBox2,
			),
		)

	//versionRight := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Width(width-w(logo)-5).Padding(1, 5, 0, 0).Render(leftBox)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		boxWrapper,
		rightLogo,
	)
}
