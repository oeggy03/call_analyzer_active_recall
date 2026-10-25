# Agent guide

This repository is the active-recall service for the Mandarin lesson ecosystem.

## Ownership

- `../call_analyzer` is the macOS desktop app. It owns SQLite and the sync
  goroutine. It pushes confirmed vocabulary and pulls it on launch.
- `../call-analyzer-lesson-db` owns the hosted Postgres migrations. Project
  ref `zoiwsytnyrdifggwleru` already has
  `20261025103511_initial_schema.sql` applied.
- This repository owns only the recall HTTP service and its Telegram delivery.

Do not edit the desktop SQLite schema or the lesson-db migrations from here.
The service reads confirmed vocabulary and does not write `study_state` or
`review_events`. Quiz replies are a later change in this repository.

## What is deployed

GCP project `call-analyzer-512714`, region `asia-southeast1`:

- Cloud Run service `active-recall`, minimum instances 0, maximum 1, no
  public invoker.
- Cloud Scheduler job `daily-vocabulary`, cron `0 9 * * *`, timezone
  `Asia/Singapore`, `POST /daily` with OIDC.
- Secret Manager names `supabase-secret`, `telegram-bot-token`, and
  `telegram-chat-id`. The Telegram bot is `@ca_chinese_active_recall_bot`.

The service URL and secret values are operational details. Read
`docs/operations.md` for the URL and the service-account names. Never copy
those secret values into git, and never reuse the desktop Mac key.

## Security and verification

- Never commit secrets, tokens, chat IDs, passwords, or connection strings
  containing passwords.
- Never log `SUPABASE_SECRET_KEY`, `TELEGRAM_BOT_TOKEN`, or
  `TELEGRAM_CHAT_ID`.
- Keep tests offline. Use `httptest` and fixtures; do not call Telegram or
  Supabase while testing.
- Do not run gcloud as part of local verification.
- Keep `/health` independent of configuration and outbound services.

Follow the contracts in the Go code and the deployment workflow in
[`docs/operations.md`](docs/operations.md).
