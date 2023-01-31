package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var palbiChatID uint64 = 312501877

const (
	monacoBilleterieURL = "https://billetterie.asmonaco.com/fr/"
	comingSoonAction    = "Prochainement"
	reserveAction       = "Réserver"
	monacoPSGMatch      = "AS MONACOPARIS SAINT GERMAIN"
)

func Scrape() {
	// Request the HTML page.
	res, err := http.Get(monacoBilleterieURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// "exists(.matchCard|text(span.teamName) == 'AS MONACOPARIS SAINT GERMAIN' && text(.matchActions) == 'Réserver')"

	// Find the review items
	doc.Find(".matchCard").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the matchTeams
		matchTeams := s.Find("span.teamName").Text()
		if matchTeams != monacoPSGMatch {
			return
		}
		actionsNode := s.Find(".matchActions")
		actionsText := strings.Trim(actionsNode.Text(), " \n")
		log.Printf("monaco psg match actions text: %s", actionsText)
		if actionsText == reserveAction {
			Notify()
		}
	})
}

func Notify() {
	log.Printf("GO RESERVE NOW!")
	err := sendTelegramMessage(palbiChatID, fmt.Sprintf("Ya se pueden reservar entradas del monaco\\! Click [aca](%s)", monacoBilleterieURL))
	if err != nil {
		log.Printf("Failed to notify that tickets are for sale!. Error: %s", err.Error())
	}
}

var client = http.Client{
	Timeout: time.Second * 10,
}

type sendMessageRequest struct {
	ChatID    uint64 `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func sendTelegramMessage(chatID uint64, text string) error {
	sendMessageReq := sendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}
	bs, err := json.Marshal(sendMessageReq)
	if err != nil {
		return fmt.Errorf("failed to marshall sendMessage request: %w", err)
	}
	log.Printf("request: %s", string(bs))
	botToken := os.Getenv("BOT_TOKEN")
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken), bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make http request to telegram: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Calling telegram api failed with status code: %d", res.StatusCode)
		defer res.Body.Close()
		bs, err := ioutil.ReadAll(res.Body)
		if err == nil {
			log.Printf("Error response: %s", string(bs))
		}
		return fmt.Errorf("sendMessage not successful")
	}
	return nil
}

var runEvery = time.Minute * 30

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting to scrape")
	sendTelegramMessage(palbiChatID, "Empezando scraper de PSG \\- AS Monaco")
	done := make(chan struct{})
	done2 := make(chan struct{})
	ticker := time.NewTicker(runEvery)

	// scrape every so often
	go func() {
		for {
			select {
			case <-done:
				log.Printf("Stopping scrape")
				return
			case <-ticker.C:
				log.Printf("scraping")
				Scrape()
			}
		}
	}()

	// cancel func
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s", sig)
		ticker.Stop()
		done <- struct{}{}
		done2 <- struct{}{}
	}()
	<-done2
}
