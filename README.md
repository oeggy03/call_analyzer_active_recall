# Mandarin Active Recall Service

This Go service reads confirmed vocabulary from the Mandarin lesson ecosystem's
hosted Supabase database and sends one morning review message to Telegram. It
currently sends the daily message only; accepting quiz replies is a later
change in this repository.

The repositories have separate responsibilities:

- `call_analyzer` is the macOS desktop app and keeps the local SQLite database.
- `call-analyzer-lesson-db` owns the shared Supabase Postgres migrations.
- This repository owns the active-recall service and its scheduled Telegram
  delivery.

The scheduled delivery is at 9:00 every morning in the `Asia/Singapore`
timezone. The service exposes:

- `GET /health` — returns `ok` without reading secrets or contacting another
  service.
- `POST /daily` — selects up to ten confirmed words and sends the message.

## Cloud Run deployment

The intended Cloud Run settings are:

- Region: `asia-southeast1`
- Unauthenticated access: disabled
- Maximum instances: `1`
- Minimum instances: `0`

Cloud Scheduler should call `POST /daily` with OIDC authentication. See
[`docs/operations.md`](docs/operations.md) for the schedule and secret names.

The macOS desktop app's Supabase key must not be reused here. This service uses
its own `SUPABASE_SECRET_KEY`, supplied through Secret Manager in production.

## Configuration

`PORT` defaults to `8080`. The daily endpoint requires:

- `SUPABASE_SECRET_KEY`
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`

`SUPABASE_URL` is optional and defaults to the hosted project URL used by the
lesson ecosystem. Copy `.env.example` for the local variable names, but never
commit populated environment files.

## Local development

Run the offline test suite with:

```sh
gofmt -w *.go
go test ./...
go vet ./...
```

The tests use `httptest`; they do not contact Supabase or Telegram.
