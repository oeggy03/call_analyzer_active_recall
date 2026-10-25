Read `AGENTS.md` for ownership, the live deployment, and offline verification.
Read [`docs/operations.md`](docs/operations.md) before changing Cloud Run,
Cloud Scheduler, or Secret Manager.

The 9:00 `Asia/Singapore` job `daily-vocabulary` is already enabled. Do not
create a second service or schedule. The desktop app and the Postgres
migrations are sibling repositories. This service only reads confirmed
vocabulary. Never commit secrets or reuse the Mac desktop Supabase key.
