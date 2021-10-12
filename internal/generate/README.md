# Generate

Generate provides fake data in a deterministic way

## Priority
1. `NULL`
2. `PG Default Val`
3. [Values from Default Values Table](#default-values)
4. [Values from Unique Values Table](#unique-values)

## Default Values

| Data Type   | Values                                                                                      |
| ----------- | ------------------------------------------------------------------------------------------- |
| integer     | 0, 1, -1, 2147483647, -2147483648                                                           |
| bool        | false, true                                                                                 |
| numeric     | 0.00, 1.23, -1.23, NaN                                                                      |
|             | + 1000___.___0001 i.e. 131072 digits before decimal point, 16383 digits after decimal point |
|             | - 1000___.___0001 i.e. 131072 digits before decimal point, 16383 digits after decimal point |
| text        | hello world, E'3?!-+@.(\u0001)ñ水불ツ😂'                                                     |
| timestamptz | 1991-11-25 12:34:56.123456+07                                                               |
|             | 4714-11-24 00:22:00+00:22 BC, 294276-12-31 23:59:59.999999+00                               |
|             | infinity, -infinity                                                                         |
| date        | 1991-11-11, 4714-11-24 BC, 5874897-12-31, infinity, -infinity                               |
| bytea       | 'hello'::bytea, 'mañana €5,90'::bytea, '\x00'                                               |
| uuid        | 00010203-0405-0607-0809-0a0b0c0d0e0f, 00000000-0000-0000-0000-000000000000                  |

## Unique Values

| Data Type   | Values                                  |
| ----------- | --------------------------------------- |
| integer     | 100 + `num`                             |
| bool        | `num` even false, `num` odd true        |
| numeric     | `num`.`num`                             |
| text        | unique_`num`                            |
| timestamptz | (2000 + `num`)-01-02 01:23:45.123456+00 |
| date        | (2000 + `num`)-01-02                    |
| bytea       | unique_`num`::bytea                     |
| uuid        | 00000000-0000-0000-0...`num`::hex       |

## Special Notes

`text` type on postgresql has limitations (1) does not accept unicode NULL character i.e. \x00, and (2) does not accept anything outside the client encoding character set (usually utf-8)

`utf-8` is represented by 1-4 bytes

- one byte: encompasses all 128 US-ASCII characters (of which 32-126 inclusive are printing)
- two bytes: encompasses next 1, 920 characters. This covers almost all Latin alphabets and Greek, Cyrillic, Hebrew, Arabic etc and Combining Diacritic Marks
- three bytes: encompasses virtually all characters in common use i.e. Chinese, Japanese, Korean (CJK)
- four bytes: encompasses less common CJK characters, math symbols, emojis