package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runBotCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(testMessageCmd)
}

var rootCmd = &cobra.Command{
	Use: "monacopsg",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
