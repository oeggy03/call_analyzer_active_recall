# Daily recall

This skill describes the live morning message. Read
[`reference.md`](reference.md) for the message format and the deployment
names. The desktop sync and the Postgres migrations are sibling
repositories. This service reads confirmed vocabulary with its own Supabase
secret. Do not reuse the Mac desktop key.

## Morning message

`POST /daily` reads confirmed vocabulary from the shared Supabase PostgREST
endpoint, shuffles the rows with an injected `math/rand/v2` random source, and
sends at most ten rows to Telegram. The first sense by lowest `sort_order` is
shown as the meaning. A traditional form is shown in parentheses only when it
is non-empty and differs from the simplified form.

An empty confirmed set still sends exactly:

```text
No confirmed words to review today.
```

## Environment

- `SUPABASE_URL` is optional and defaults to the hosted lesson project URL.
- `SUPABASE_SECRET_KEY` is required for `/daily`.
- `TELEGRAM_BOT_TOKEN` is required for `/daily`.
- `TELEGRAM_CHAT_ID` is required for `/daily`.
- `PORT` defaults to `8080`.

The secret, token, and chat ID must never be committed or logged. Tests use
local `httptest` servers and do not call external services.

## Scope

Cloud Scheduler job `daily-vocabulary` already calls `POST /daily` at 9:00
`Asia/Singapore` on Cloud Run service `active-recall` in project
`call-analyzer-512714`. Do not create a second service or job.

This version sends the morning message only. Quiz replies and writes to
`study_state` are a later change in this repository.
