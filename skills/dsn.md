# DSN Generator Skill

Use the DSN generator script to produce a PostgreSQL DSN from a `.env` file.

## Command
- `python3 scripts/py/dsn_generator/main.py --env <path-to-env>`

## Expected output (with standard dev env)
- `postgres://devuser:devpass@localhost:5433/devdb?sslmode=disable`

## Notes for LLM agents
- Requires Python 3 (use `python3`).
- `.env` keys required: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`.
- Values may be quoted or unquoted; script strips surrounding quotes.
- Percent-encoding is applied to user/password, so credentials with special characters are safe.
- Script exits with an error if keys are missing or the env file is not found.
