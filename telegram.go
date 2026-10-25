package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIBaseURL = "https://api.telegram.org"

type TelegramClient struct {
	Token      string
	ChatID     string
	Endpoint   string
	HTTPClient *http.Client
}

func (c *TelegramClient) Send(ctx context.Context, text string) error {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = telegramAPIBaseURL + "/bot" + c.Token + "/sendMessage"
	}

	payload, err := json.Marshal(struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: c.ChatID,
		Text:   text,
	})
	if err != nil {
		return redactedError(fmt.Sprintf("encoding Telegram request: %v", err), c.Token)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return redactedError(fmt.Sprintf("creating Telegram request: %v", err), c.Token)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return redactedError(fmt.Sprintf("Telegram request failed: %v", err), c.Token)
	}
	defer response.Body.Close()

	body, err := readResponseBody(response.Body, c.Token)
	if err != nil {
		return redactedError(fmt.Sprintf("reading Telegram response: %v", err), c.Token)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return redactedError(
			fmt.Sprintf("Telegram returned %s: %s", response.Status, string(body)),
			c.Token,
		)
	}
	return nil
}
