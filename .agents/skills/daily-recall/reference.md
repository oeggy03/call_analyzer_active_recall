# Daily recall reference

## Vocabulary read

The service requests confirmed entries from:

```text
/rest/v1/vocabulary_entries?status=eq.confirmed&select=id,simplified,traditional,reading,senses(gloss,sort_order)&limit=1000
```

The Supabase secret is sent in both the `apikey` header and the
`Authorization: Bearer` header. Non-2xx responses are failures, and response
errors redact the secret.

## Selection and formatting

1. Copy the fetched entries.
2. Shuffle the copy with the selector's injected `math/rand/v2` `*rand.Rand`.
3. Keep the first ten entries, or all entries when fewer than ten exist.
4. Use the sense with the lowest `sort_order` as the meaning.
5. Send plain text with no Telegram parse mode:

```text
Today's words

1. 简体（繁體）
reading
meaning
```

Words are numbered from one and separated by a blank line. The parenthetical
traditional form is omitted when it is empty or matches the simplified form.

## Delivery configuration

`SUPABASE_URL` defaults to the hosted project URL and can be overridden.
`SUPABASE_SECRET_KEY`, `TELEGRAM_BOT_TOKEN`, and `TELEGRAM_CHAT_ID` are
required for `/daily`. The bot token is used in the Telegram
`sendMessage` endpoint and is redacted from failure text.

`GET /health` returns `ok` without reading any of these values or contacting
another service.

## Deployment

The running names are GCP project `call-analyzer-512714`, region
`asia-southeast1`, Cloud Run service `active-recall`, and Scheduler job
`daily-vocabulary`. Secret Manager holds `supabase-secret`,
`telegram-bot-token`, and `telegram-chat-id`. The bot username is
`@ca_chinese_active_recall_bot`. The service URL and account emails are in
`docs/operations.md`. Do not copy secret values into this repository.

## Future scope

Quiz replies are intentionally not implemented yet. A future change may add
reply handling and `study_state` writes, but this service currently reads
confirmed vocabulary only.
