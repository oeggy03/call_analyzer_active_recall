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

## Scheduled delivery

Create a Cloud Scheduler job with:

- Cron: `0 9 * * *`
- Timezone: `Asia/Singapore`
- Method: `POST`
- URI: `<CLOUD_RUN_SERVICE_URL>/daily`
- Authentication: OIDC using a service account authorized to invoke the
  Cloud Run service

The Cloud Run service should remain private, with no unauthenticated access.
Use region `asia-southeast1`, maximum instances `1`, and no minimum instances.

## Health check

`GET <CLOUD_RUN_SERVICE_URL>/health` returns `ok`. It does not read Secret
Manager values or call Supabase or Telegram.
