#!/usr/bin/env python3
"""Generate a PostgreSQL DSN string from a .env file."""

from __future__ import annotations

import argparse
from pathlib import Path
from typing import Dict
from urllib.parse import quote_plus

REQUIRED_KEYS = (
    "DB_HOST",
    "DB_PORT",
    "DB_USER",
    "DB_PASSWORD",
    "DB_NAME",
    "DB_SSLMODE",
)


def parse_env_file(path: Path) -> Dict[str, str]:
    data: Dict[str, str] = {}
    for raw_line in path.read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("export "):
            line = line[len("export ") :].strip()
        if "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip()
        if (value.startswith("\"") and value.endswith("\"")) or (
            value.startswith("'") and value.endswith("'")
        ):
            value = value[1:-1]
        data[key] = value
    return data


def build_postgres_dsn(values: Dict[str, str]) -> str:
    missing = [key for key in REQUIRED_KEYS if not values.get(key)]
    if missing:
        joined = ", ".join(missing)
        raise ValueError(f"Missing required keys: {joined}")

    user = quote_plus(values["DB_USER"])
    password = quote_plus(values["DB_PASSWORD"])
    host = values["DB_HOST"]
    port = values["DB_PORT"]
    name = values["DB_NAME"]
    sslmode = values["DB_SSLMODE"]

    return f"postgres://{user}:{password}@{host}:{port}/{name}?sslmode={sslmode}"


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Generate a PostgreSQL DSN string from a .env file."
    )
    parser.add_argument(
        "--env",
        default=".env",
        help="Path to the .env file (default: .env)",
    )
    args = parser.parse_args()

    env_path = Path(args.env)
    if not env_path.exists():
        parser.error(f".env file not found: {env_path}")

    values = parse_env_file(env_path)
    try:
        dsn = build_postgres_dsn(values)
    except ValueError as exc:
        parser.error(str(exc))

    print(dsn)


if __name__ == "__main__":
    main()
