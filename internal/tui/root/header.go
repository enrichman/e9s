package root

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/enrichman/e9s/internal/constants"
	"github.com/enrichman/e9s/internal/tui/style"
)

type Header struct {
	version string
}

func NewHeader(version string) Header {
	return Header{version: version}
}

func viewHeader(header Header) string {
	w := lipgloss.Width

	logo := logo()

	infoBoxWrapper := lipgloss.NewStyle().
		Padding(1).
		Width(width - w(logo)).
		Render(viewInfoBox(header.version))

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		infoBoxWrapper,
		logo,
	)
}

func viewInfoBox(version string) string {
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
