# Daily recall

This skill describes the first version of the active-recall service. Read
[`reference.md`](reference.md) for the detailed message and integration
contract.

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

This version sends the morning message only. Quiz replies and writes to
`study_state` are a later change in this repository.
