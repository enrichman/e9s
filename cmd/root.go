package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enrichman/e9s/internal/tui/root"
	"github.com/spf13/cobra"
)

const AppName = `___________________       
\_   _____/   __   \______
 |    __)_\____    /  ___/
 |        \  /    /\___ \ 
/_______  / /____//____  >
        \/             \/ 
`

var (
	Version = "0.0.0-dev"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   AppName,
		Short: `E9s is a CLI to view and manage your Epinio cluster.`,
		RunE:  rootRunE,
	}

	rootCmd.AddCommand(
		newVersionCmd(),
	)

	return rootCmd
}

func rootRunE(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(
		root.NewRootModel(Version),
		tea.WithAltScreen(),
	)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = p.Run()
	return err
}
