package cmd

import (
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enrichman/e9s/internal/tui/root"
	"github.com/enrichman/e9s/pkg/client"
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
	insecureTransport := http.DefaultTransport.(*http.Transport).Clone()
	insecureTransport.TLSClientConfig.InsecureSkipVerify = true

	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: insecureTransport,
	}

	ep := client.NewClient(c, "https://epinio.172.21.0.4.omg.howdoi.website/api/v1")
	ep.Username = "admin"
	ep.Password = "password"

	p := tea.NewProgram(
		root.NewRootModel(ep, Version),
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
