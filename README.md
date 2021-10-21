# PGMock

PGMock generates data to provide a standard base state for testing migrations.

Given a schema, generate mock data in a way that

- enumerates all column combinations
- is idempotent
- is deterministic

Note: currently in development with no guarantees for compatabilty between commits

## Database Connection

Uses [PGConnect](https://www.github.com/justasable/pgconnect)

## Testing

Requires `psql`, `dropdb` and `createdb`
