package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thepalbi/monaco-psg-bot/pkg"
)

var userID, text string

var testMessageCmd = &cobra.Command{
	Use:   "test-message",
	Short: "Send a test message in telegram",
	Run: func(cmd *cobra.Command, args []string) {
		parsedUserID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("failed to parse user ID: %s", err)
			return
		}
		pkg.SendTelegramMessage(botToken, uint64(parsedUserID), text)
	},
}

func init() {
	testMessageCmd.Flags().StringVar(&userID, "user-id", "", "User id to send message to")
	testMessageCmd.Flags().StringVar(&botToken, "bot-token", "", "Bot token")
	testMessageCmd.Flags().StringVar(&text, "text", "holis", "Text to send")
	testMessageCmd.MarkFlagRequired("user-id")
	testMessageCmd.MarkFlagRequired("bot-token")
}
