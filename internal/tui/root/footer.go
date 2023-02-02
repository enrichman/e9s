package root

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/tui/style"
)

type Footer struct {
	currentPage string
	action      string
}

func NewFooter(initialPage string) Footer {
	return Footer{currentPage: initialPage}
}

func updateFooter(footer Footer, currentPage, action string) Footer {
	footer.currentPage = currentPage
	footer.action = action
	return footer
}

func viewFooter(footer Footer) string {
	views := []string{}

	currentPage := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Bold(true).
		Background(style.ColorOrange).
		Render(footer.currentPage)

	views = append(views, currentPage)

	if footer.action != "" {
		currentAction := lipgloss.NewStyle().
			Margin(0, 1).
			Padding(0, 1).
			Bold(true).
			Background(style.ColorOrange).
			Render(footer.action)

		views = append(views, currentAction)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, views...)
}
