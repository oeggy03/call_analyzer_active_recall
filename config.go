package main

import (
	"fmt"
	"os"
	"strings"
)

const defaultSupabaseURL = "https://zoiwsytnyrdifggwleru.supabase.co"

type Config struct {
	SupabaseURL    string
	SupabaseSecret string
	TelegramToken  string
	TelegramChatID string
}

func loadConfig() (Config, error) {
	return configFromEnv(os.Getenv)
}

func configFromEnv(getenv func(string) string) (Config, error) {
	supabaseURL := strings.TrimSpace(getenv("SUPABASE_URL"))
	if supabaseURL == "" {
		supabaseURL = defaultSupabaseURL
	} else {
		supabaseURL = strings.TrimRight(supabaseURL, "/")
	}

	config := Config{
		SupabaseURL:    supabaseURL,
		SupabaseSecret: getenv("SUPABASE_SECRET_KEY"),
		TelegramToken:  getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID: getenv("TELEGRAM_CHAT_ID"),
	}

	for _, required := range []struct {
		name  string
		value string
	}{
		{name: "SUPABASE_SECRET_KEY", value: config.SupabaseSecret},
		{name: "TELEGRAM_BOT_TOKEN", value: config.TelegramToken},
		{name: "TELEGRAM_CHAT_ID", value: config.TelegramChatID},
	} {
		if strings.TrimSpace(required.value) == "" {
			return Config{}, fmt.Errorf("missing required environment variable %s", required.name)
		}
	}

	return config, nil
}
