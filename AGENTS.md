# Agent guide

This repository is the active-recall service for the Mandarin lesson ecosystem.

## Ownership

- `call_analyzer` is the macOS desktop app and owns its local SQLite database.
- `call-analyzer-lesson-db` owns the hosted Supabase Postgres migrations.
- This repository owns only the recall HTTP service and its Telegram delivery.

Do not edit the desktop SQLite schema or the `call-analyzer-lesson-db`
migrations from this repository. The service reads confirmed vocabulary rows
from Supabase and does not write `study_state` yet. Quiz replies are a later
change in this repository.

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
