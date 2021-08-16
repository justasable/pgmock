# PGMock

PGMock generates data to provide a standard base state for testing migrations.

Given a schema, generate mock data in a way that

- enumerates all column combinations
- is idempotent
- is deterministic

Note: currently in development with no guarantees for compatabilty between commits

### Supported Data Type

| Data Type | Behavior |
| --------- | -------- |
|           |          |

### Supported Constraints

| Constraints                 | Behavior                             |
| --------------------------- | ------------------------------------ |
| nullable                    | null, default (see PGMock Defaults)  |
| default (user overwritable) | pgmock default (see PGMock Defaults) |
| default (db generated)      | ignore                               |
| generated                   | ignore                               |
| check                       | ignore                               |

### PGMock Defaults

| Default Type | Value |
| ------------ | ----- |
|              |       |

## Roadmap

**Numeric**: integer, smallint, bigint, decimal, numeric, serial, smallserial, bigserial

**Text**: text, char(n), varchar(n)

**Other**: primary key, foreign key

## Testing

- requires `psql`
- uses `pgmock_test_db` for testing
- uses [pgx](https://github.com/jackc/pgx) defaults to connect to database (except the above database)
