# Operations

## Secret Manager

Production deployment supplies these environment variables from Secret
Manager. The repository contains names only; never put their values in Git,
logs, documentation, or request examples.

| Environment variable | Secret Manager secret |
| --- | --- |
| `SUPABASE_SECRET_KEY` | `supabase-secret` |
| `TELEGRAM_BOT_TOKEN` | `telegram-bot-token` |
| `TELEGRAM_CHAT_ID` | `telegram-chat-id` |

`SUPABASE_URL` may be configured as a non-secret environment variable. If it
is omitted, the service uses its built-in hosted project URL.

## Current deployment

GCP project `call-analyzer-512714`, region `asia-southeast1`.

- Cloud Run service `active-recall`. URL
  `https://active-recall-ac32aneera-as.a.run.app`. Minimum instances 0,
  maximum instances 1, unauthenticated access disabled. Runtime service
  account `active-recall-runner@call-analyzer-512714.iam.gserviceaccount.com`.
- Cloud Scheduler job `daily-vocabulary`. Cron `0 9 * * *`, timezone
  `Asia/Singapore`, `POST /daily` on that URL. OIDC service account
  `active-recall-invoker@call-analyzer-512714.iam.gserviceaccount.com`.
- Telegram bot `@ca_chinese_active_recall_bot`.

Do not create a second service or schedule for the morning message. Change
the existing job when the time or target changes. A manual check is
`gcloud scheduler jobs run daily-vocabulary --location=asia-southeast1`.

## Health check

`GET /health` returns `ok` and does not read secrets or call Supabase or
Telegram. The deployed service rejects unauthenticated calls, so a check from
outside GCP needs an identity token.
