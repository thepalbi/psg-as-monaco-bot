package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 10,
}

type sendMessageRequest struct {
	ChatID    uint64 `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func SendTelegramMessage(botToken string, chatID uint64, text string) error {
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
