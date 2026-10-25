package main

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
)

type vocabularySource interface {
	FetchConfirmed(context.Context) ([]VocabularyEntry, error)
}

type messageSender interface {
	Send(context.Context, string) error
}

type dailyService struct {
	source   vocabularySource
	sender   messageSender
	selector *selector
}

func newDailyService(config Config, httpClient *http.Client, random *rand.Rand) *dailyService {
	return &dailyService{
		source: &SupabaseClient{
			BaseURL:    config.SupabaseURL,
			Secret:     config.SupabaseSecret,
			HTTPClient: httpClient,
		},
		sender: &TelegramClient{
			Token:      config.TelegramToken,
			ChatID:     config.TelegramChatID,
			HTTPClient: httpClient,
		},
		selector: newSelector(random),
	}
}

func (s *dailyService) SendDaily(ctx context.Context) (int, error) {
	entries, err := s.source.FetchConfirmed(ctx)
	if err != nil {
		return 0, err
	}

	selected := s.selector.Select(entries, 10)
	if err := s.sender.Send(ctx, formatDailyMessage(selected)); err != nil {
		return 0, err
	}
	return len(selected), nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func dailyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	config, err := loadConfig()
	if err != nil {
		writeDailyFailure(w, err, nil)
		return
	}

	service := newDailyService(config, http.DefaultClient, newRandom())
	serveDaily(w, r, service, config.SupabaseSecret, config.TelegramToken)
}

func newDailyHandler(service *dailyService, secrets ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveDaily(w, r, service, secrets...)
	}
}

func serveDaily(w http.ResponseWriter, r *http.Request, service *dailyService, secrets ...string) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	count, err := service.SendDaily(r.Context())
	if err != nil {
		writeDailyFailure(w, err, secrets)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Sent  bool `json:"sent"`
		Count int  `json:"count"`
	}{
		Sent:  true,
		Count: count,
	})
}

func writeDailyFailure(w http.ResponseWriter, err error, secrets []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(struct {
		Sent  bool   `json:"sent"`
		Error string `json:"error"`
	}{
		Sent:  false,
		Error: redactSecrets(err.Error(), secrets...),
	})
}
