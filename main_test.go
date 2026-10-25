package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDailySendsEmptyDay(t *testing.T) {
	const (
		supabaseSecret = "supabase-test-secret"
		telegramToken  = "telegram-test-token"
	)

	supabase := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/v1/vocabulary_entries" {
			t.Errorf("unexpected Supabase path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("status"); got != "eq.confirmed" {
			t.Errorf("unexpected status filter: %s", got)
		}
		if got := r.URL.Query().Get("select"); got != "id,simplified,traditional,reading,senses(gloss,sort_order)" {
			t.Errorf("unexpected select: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1000" {
			t.Errorf("unexpected limit: %s", got)
		}
		if got := r.Header.Get("apikey"); got != supabaseSecret {
			t.Errorf("unexpected apikey header")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+supabaseSecret {
			t.Errorf("unexpected authorization header")
		}
		_, _ = w.Write([]byte("[]"))
	}))
	defer supabase.Close()

	var telegramCalls int
	var telegramText string
	telegram := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegramCalls++
		var payload struct {
			ChatID string `json:"chat_id"`
			Text   string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode Telegram request: %v", err)
		}
		if payload.ChatID != "chat-test-id" {
			t.Errorf("unexpected chat id")
		}
		telegramText = payload.Text
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer telegram.Close()

	service := &dailyService{
		source: &SupabaseClient{
			BaseURL:    supabase.URL,
			Secret:     supabaseSecret,
			HTTPClient: telegram.Client(),
		},
		sender: &TelegramClient{
			Token:      telegramToken,
			ChatID:     "chat-test-id",
			Endpoint:   telegram.URL,
			HTTPClient: telegram.Client(),
		},
		selector: newSelector(rand.New(rand.NewPCG(1, 2))),
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/daily", nil)
	newDailyHandler(service, supabaseSecret, telegramToken).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	var response struct {
		Sent  bool `json:"sent"`
		Count int  `json:"count"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode daily response: %v", err)
	}
	if !response.Sent || response.Count != 0 {
		t.Fatalf("unexpected daily response: sent=%v count=%d", response.Sent, response.Count)
	}
	if telegramCalls != 1 {
		t.Fatalf("Telegram calls: got %d, want 1", telegramCalls)
	}
	if telegramText != "No confirmed words to review today." {
		t.Fatalf("unexpected empty-day message")
	}
}

func TestDailySelectsTenAndSameSeedIsStable(t *testing.T) {
	entries := testEntries(12)
	supabase := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			t.Errorf("encode vocabulary: %v", err)
		}
	}))
	defer supabase.Close()

	var telegramMessages []string
	telegram := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode Telegram request: %v", err)
		}
		telegramMessages = append(telegramMessages, payload.Text)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer telegram.Close()

	const (
		supabaseSecret = "supabase-test-secret"
		telegramToken  = "telegram-test-token"
	)
	newService := func() *dailyService {
		return &dailyService{
			source: &SupabaseClient{
				BaseURL:    supabase.URL,
				Secret:     supabaseSecret,
				HTTPClient: telegram.Client(),
			},
			sender: &TelegramClient{
				Token:      telegramToken,
				ChatID:     "chat-test-id",
				Endpoint:   telegram.URL,
				HTTPClient: telegram.Client(),
			},
			selector: newSelector(rand.New(rand.NewPCG(42, 99))),
		}
	}

	for range 2 {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/daily", nil)
		newDailyHandler(newService(), supabaseSecret, telegramToken).ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("unexpected status: %d", recorder.Code)
		}
		var response struct {
			Sent  bool `json:"sent"`
			Count int  `json:"count"`
		}
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("decode daily response: %v", err)
		}
		if !response.Sent || response.Count != 10 {
			t.Fatalf("unexpected daily response: sent=%v count=%d", response.Sent, response.Count)
		}
	}

	if len(telegramMessages) != 2 {
		t.Fatalf("Telegram messages: got %d, want 2", len(telegramMessages))
	}
	if telegramMessages[0] != telegramMessages[1] {
		t.Fatalf("same seed produced different messages")
	}
	if strings.Count(telegramMessages[0], "\n\n") != 10 {
		t.Fatalf("message did not contain exactly 10 separated words")
	}
}

func TestFormatDailyMessageOmitsMatchingTraditional(t *testing.T) {
	message := formatDailyMessage([]VocabularyEntry{{
		Simplified:  "你好",
		Traditional: "你好",
		Reading:     "ni3 hao3",
		Senses:      []Sense{{Gloss: "hello", SortOrder: 0}},
	}})

	want := "Today's words\n\n1. 你好\nnǐ hǎo\nhello"
	if message != want {
		t.Fatalf("unexpected message")
	}
}

func TestFormatDailyMessageRendersToneMarks(t *testing.T) {
	message := formatDailyMessage([]VocabularyEntry{
		{
			Simplified: "黑五",
			Reading:    "hei1 wu3",
			Senses:     []Sense{{Gloss: "Black Friday. After it, Q1 traffic growth is manageable.", SortOrder: 0}},
		},
		{
			Simplified: "利用率",
			Reading:    "li4 yong4 lü4",
			Senses:     []Sense{{Gloss: "utilization. IES 中台's lot is about 381 cards.", SortOrder: 0}},
		},
	})

	want := "Today's words\n\n1. 黑五\nhēi wǔ\nBlack Friday. After it, Q1 traffic growth is manageable.\n\n2. 利用率\nlì yòng lǜ\nutilization. IES 中台's lot is about 381 cards."
	if message != want {
		t.Fatalf("unexpected message:\n%s", message)
	}
}

func TestFormatDailyMessageUsesLowestSortOrderSense(t *testing.T) {
	message := formatDailyMessage([]VocabularyEntry{{
		Simplified: "算法",
		Reading:    "suan4 fa3",
		Senses: []Sense{
			{Gloss: "later meaning", SortOrder: 10},
			{Gloss: "algorithm", SortOrder: 1},
			{Gloss: "first meaning", SortOrder: -1},
		},
	}})

	if !strings.HasSuffix(message, "first meaning") {
		t.Fatalf("message did not use the lowest sort_order sense")
	}
	if strings.Contains(message, "later meaning") || strings.Contains(message, "algorithm") {
		t.Fatalf("message included a non-lowest sense")
	}
}

func TestSupabaseFailureDoesNotCallTelegramOrLeakSecret(t *testing.T) {
	const (
		supabaseSecret = "supabase-test-secret"
		telegramToken  = "telegram-test-token"
	)

	supabase := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "upstream failure: %s", supabaseSecret)
	}))
	defer supabase.Close()

	var telegramCalls int
	telegram := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegramCalls++
	}))
	defer telegram.Close()

	service := &dailyService{
		source: &SupabaseClient{
			BaseURL:    supabase.URL,
			Secret:     supabaseSecret,
			HTTPClient: telegram.Client(),
		},
		sender: &TelegramClient{
			Token:      telegramToken,
			ChatID:     "chat-test-id",
			Endpoint:   telegram.URL,
			HTTPClient: telegram.Client(),
		},
		selector: newSelector(rand.New(rand.NewPCG(1, 2))),
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/daily", nil)
	newDailyHandler(service, supabaseSecret, telegramToken).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if telegramCalls != 0 {
		t.Fatalf("Telegram calls: got %d, want 0", telegramCalls)
	}
	if strings.Contains(recorder.Body.String(), supabaseSecret) {
		t.Fatal("Supabase secret appeared in HTTP response")
	}
}

func TestTelegramFailureDoesNotLeakBotToken(t *testing.T) {
	const (
		supabaseSecret = "supabase-test-secret"
		telegramToken  = "telegram-test-token"
	)

	supabase := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(testEntries(1))
	}))
	defer supabase.Close()

	telegram := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "bad bot token: %s", telegramToken)
	}))
	defer telegram.Close()

	service := &dailyService{
		source: &SupabaseClient{
			BaseURL:    supabase.URL,
			Secret:     supabaseSecret,
			HTTPClient: telegram.Client(),
		},
		sender: &TelegramClient{
			Token:      telegramToken,
			ChatID:     "chat-test-id",
			Endpoint:   telegram.URL,
			HTTPClient: telegram.Client(),
		},
		selector: newSelector(rand.New(rand.NewPCG(1, 2))),
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/daily", nil)
	newDailyHandler(service, supabaseSecret, telegramToken).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), telegramToken) {
		t.Fatal("Telegram bot token appeared in HTTP response")
	}
}

func TestHealthWorksWithoutDailySecrets(t *testing.T) {
	t.Setenv("SUPABASE_SECRET_KEY", "")
	t.Setenv("TELEGRAM_BOT_TOKEN", "")
	t.Setenv("TELEGRAM_CHAT_ID", "")

	healthRecorder := httptest.NewRecorder()
	healthRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthHandler(healthRecorder, healthRequest)
	if healthRecorder.Code != http.StatusOK || healthRecorder.Body.String() != "ok" {
		t.Fatalf("health response was not ok")
	}

	dailyRecorder := httptest.NewRecorder()
	dailyRequest := httptest.NewRequest(http.MethodPost, "/daily", nil)
	dailyHandler(dailyRecorder, dailyRequest)
	if dailyRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("missing configuration status: got %d", dailyRecorder.Code)
	}
}

func testEntries(count int) []VocabularyEntry {
	entries := make([]VocabularyEntry, count)
	for index := range entries {
		entries[index] = VocabularyEntry{
			ID:          fmt.Sprintf("entry-%d", index),
			Simplified:  fmt.Sprintf("词%d", index),
			Traditional: fmt.Sprintf("詞%d", index),
			Reading:     fmt.Sprintf("ci2 %d", index),
			Senses:      []Sense{{Gloss: fmt.Sprintf("word %d", index), SortOrder: 0}},
		}
	}
	return entries
}
