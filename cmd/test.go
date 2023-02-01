package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thepalbi/monaco-psg-bot/pkg"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the scraping",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Scrape(func() {
			fmt.Println("Success!")
		})
	},
}
