# Generate

Generate provides fake data in a deterministic way

## Values

| Data Type   | Values                                                                                                            |
| ----------- | ----------------------------------------------------------------------------------------------------------------- |
| integer     | 0, 1, -1, 2147483648, -2147483648                                                                                 |
| numeric     | 0.00, -1.23, 1.23, NaN                                                                                            |
|             | ± 1000___.___0001 i.e. 131072 digits before decimal point, 16383 digits after decimal point                       |
| text        | hello world, E'3?!-+@.(\u0001)ñ水불ツ😂'                                                                           |
| timestamptz | 2021-11-01 12:34:56.123456+03, -infinity, infinity                                                                |
| date        | 1991-11-11, -infinity, infinity                                                                                   |
| bytea       | 'hello'::bytea, 'mañana €5,90'::bytea, '\x00'                                                                     |
| uuid        | 00010203-0405-0607-0809-0a0b0c0d0e0f, 00000000-0000-0000-0000-000000000000                                        |

## Special Notes

`text` type on postgresql has limitations (1) does not accept unicode NULL character i.e. \x00, and (2) does not accept anything outside the client encoding character set (usually utf-8)

`utf-8` is represented by 1-4 bytes

- one byte: encompasses all 128 US-ASCII characters (of which 32-126 inclusive are printing)
- two bytes: encompasses next 1, 920 characters. This covers almost all Latin alphabets and Greek, Cyrillic, Hebrew, Arabic etc and Combining Diacritic Marks
- three bytes: encompasses virtually all characters in common use i.e. Chinese, Japanese, Korean (CJK)
- four bytes: encompasses less common CJK characters, math symbols, emojis
