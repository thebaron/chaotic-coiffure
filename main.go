package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/thebaron/chaotic-coiffure/pkg/config"
	"github.com/thebaron/chaotic-coiffure/pkg/view"
)

func main() {
	var configPath string

	rootCmd := &cobra.Command{
		Use:   "chaotic-coiffure",
		Short: "An AI enabled subsonic playlist creator",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := config.LoadConfig(configPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			p := tea.NewProgram(view.InitialModel(c))
			_, err = p.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "path to the config.yaml file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
