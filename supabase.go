package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SupabaseClient struct {
	BaseURL    string
	Secret     string
	HTTPClient *http.Client
}

func (c *SupabaseClient) FetchConfirmed(ctx context.Context) ([]VocabularyEntry, error) {
	endpoint, err := url.Parse(strings.TrimRight(c.BaseURL, "/") + "/rest/v1/vocabulary_entries")
	if err != nil {
		return nil, redactedError(fmt.Sprintf("invalid Supabase URL: %v", err), c.Secret)
	}

	query := endpoint.Query()
	query.Set("status", "eq.confirmed")
	query.Set("select", "id,simplified,traditional,reading,senses(gloss,sort_order)")
	query.Set("limit", "1000")
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, redactedError(fmt.Sprintf("creating Supabase request: %v", err), c.Secret)
	}
	request.Header.Set("apikey", c.Secret)
	request.Header.Set("Authorization", "Bearer "+c.Secret)
	request.Header.Set("Accept", "application/json")

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, redactedError(fmt.Sprintf("Supabase request failed: %v", err), c.Secret)
	}
	defer response.Body.Close()

	body, err := readResponseBody(response.Body, c.Secret)
	if err != nil {
		return nil, redactedError(fmt.Sprintf("reading Supabase response: %v", err), c.Secret)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, redactedError(
			fmt.Sprintf("Supabase returned %s: %s", response.Status, string(body)),
			c.Secret,
		)
	}

	var entries []VocabularyEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, redactedError(fmt.Sprintf("decoding Supabase response: %v", err), c.Secret)
	}
	return entries, nil
}
