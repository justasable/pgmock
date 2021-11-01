# PGMock

PGMock generates data to provide a standard base state for testing migrations.

Given a schema, generate mock data in a way that

- enumerates all column combinations
- is idempotent
- is deterministic

Note: currently in development with no guarantees for compatabilty between commits

Compatible wiht PG14 onwards

## Database Connection

Uses [PGConnect](https://www.github.com/justasable/pgconnect)

`PGMOCK_SETUP_SCRIPT` must be set for pgmock to rebuild database to known state

## Data Generation Policy
1. `NULL` value
2. [Values from Test Values Table](#test-values)
3. database provided `DEFAULT` value (except boolean)
4. [Values from Unique Values Table](#unique-values)

The reason for this order is because the database default value can clash with our test values causing an error and preventing other test values from being inserted. Hence this order delays any potential errors up to the inevitable moment, creating a greater chance for successful row generation.

## Test Values

| Data Type     | Values                                                                                      |
| -----------   | ------------------------------------------------------------------------------------------- |
| integer       | 0, 1, -1, 2147483647, -2147483648                                                           |
| bool          | true, false                                                                                 |
| numeric       | 0.00, 1.23, -1.23, NaN, Infinity, -Infinity                                                 |
|               | + 1000___.___0001 i.e. 131072 digits before decimal point, 16383 digits after decimal point |
|               | - 1000___.___0001 i.e. 131072 digits before decimal point, 16383 digits after decimal point |
| numeric(p, s) | 0.00, 1.23, -1.23, NaN, (max val), (min val)                                                |
| text          | hello world, E'3?!-+@.(\u0001)ñ水불ツ😂'                                                      |
| timestamptz   | 1991-11-25 12:34:56.123456+07                                                               |
|               | 4714-11-24 00:22:00+00:22 BC, 294276-12-31 23:59:59.999999+00                               |
|               | infinity, -infinity                                                                         |
| date          | 1991-11-11, 4714-11-24 BC, 5874897-12-31, infinity, -infinity                               |
| bytea         | 'hello'::bytea, 'mañana €5,90'::bytea, '\x00'                                               |
| uuid          | 00010203-0405-0607-0809-0a0b0c0d0e0f, 00000000-0000-0000-0000-000000000000                  |

| Constraint      | Values                           |
| --------------- | -------------------------------- |
| GENERATED       | (db default)                     |
| IDENTITY ALWAYS | (db default)                     |
| FOREIGN KEY     | (value from first generated row) |

## Unique Values

| Data Type     | Values                                         |
| ------------- | ---------------------------------------------- |
| integer       | 100, 101, 102...                               |
| bool          | nil                                            |
| numeric       | 1.1, 10.01, 100.001...                         |
| numeric(5, 2) | 000.01, 000.02, 000.03...                      |
| text          | unique_0, unique_1, unique_2...                |
| timestamptz   | (2000 + `0, 1, 2...`)-01-02 01:23:45.123456+00 |
| date          | (2000 + `0, 1, 2...`)-01-02                    |
| bytea         | unique_`0, 1, 2...`::bytea                     |
| uuid          | 00000000-0000-0000-0...`1, 2, 3...`::hex       |

## Special Notes

`text` type on postgresql has limitations (1) does not accept unicode NULL character i.e. \x00, and (2) does not accept anything outside the client encoding character set (usually utf-8)

`utf-8` is represented by 1-4 bytes

- one byte: encompasses all 128 US-ASCII characters (of which 32-126 inclusive are printing)
- two bytes: encompasses next 1, 920 characters. This covers almost all Latin alphabets and Greek, Cyrillic, Hebrew, Arabic etc and Combining Diacritic Marks
- three bytes: encompasses virtually all characters in common use i.e. Chinese, Japanese, Korean (CJK)
- four bytes: encompasses less common CJK characters, math symbols, emojis


## Testing

Requires `psql`, `dropdb` and `createdb`
