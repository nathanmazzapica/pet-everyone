# DSN Generator

Generate a PostgreSQL DSN string from a `.env` file.

## Usage

```
python3 skills/dsn_generator/main.py --env .env
```

Expected output for the provided configuration:

```
postgres://devuser:devpass@localhost:5433/devdb?sslmode=disable
```
