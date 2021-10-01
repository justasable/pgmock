# PGMock

PGMock generates data to provide a standard base state for testing migrations.

Given a schema, generate mock data in a way that

- enumerates all column combinations
- is idempotent
- is deterministic

Note: currently in development with no guarantees for compatabilty between commits

## Connection

Database connection params must be specified explicitly for safety

- `PGMOCK_HOST`
- `PGMOCK_PORT`
- `PGMOCK_DATABASE`
- `PGMOCK_USER`
- `PGMOCK_PASSWORD` (optional)
- `PGMOCK_SETUP_SCRIPT_PATH` (optional)

## Testing

Requires `psql`, `dropdb` and `createdb`
