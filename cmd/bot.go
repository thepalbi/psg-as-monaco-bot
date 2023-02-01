package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/thepalbi/monaco-psg-bot/pkg"
)

const (
	palbiChatID uint64 = 312501877
	runEvery           = time.Minute * 30
)

var botToken string
var userIDs []int64 = []int64{}

func init() {
	runBotCmd.Flags().StringVar(&botToken, "bot-token", "", "Telegram bot token")
	runBotCmd.Flags().Int64SliceVar(&userIDs, "user-ids", nil, "User ids to notify if success")
	runBotCmd.MarkFlagRequired("bot-token")
}

var runBotCmd = &cobra.Command{
	Use:   "run-bot",
	Short: "Start scraper bot",
	Run: func(cmd *cobra.Command, args []string) {
		runBot(func(message string) {
			for _, uid := range userIDs {
				err := pkg.SendTelegramMessage(botToken, uint64(uid), message)
				if err != nil {
					log.Printf("Failed to notify that tickets are for sale!. Error: %s", err.Error())
				}
			}
		})
	},
}

type notfier = func(string)

func runBot(notifyAll notfier) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting to scrape")
	notifyAll("Empezando scraper de PSG \\- AS Monaco")
	done := make(chan struct{})
	done2 := make(chan struct{})
	ticker := time.NewTicker(runEvery)

	// scrape every so often
	go func() {
		for {
			// do one scrape first
			pkg.Scrape(func() {
				log.Printf("GO RESERVE NOW!")
				notifyAll(fmt.Sprintf("Ya se pueden reservar entradas del monaco\\! Click [aca](%s)", pkg.MonacoBilleterieURL))
			})
			select {
			case <-done:
				log.Printf("Stopping scrape")
				return
			case <-ticker.C:
			}
		}
	}()

	// cancel func
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s", sig)
		notifyAll("Por alguna razón se apago el bot")
		ticker.Stop()
		done <- struct{}{}
		done2 <- struct{}{}
	}()
	<-done2
}
