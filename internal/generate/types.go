package generate

import (
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgtype"
)

func IntegerDefaults() []int {
	return []int{0, 1, -1, 2147483648, -2147483648}
}

func IntegerUnique(num int) int {
	return 100 + num
}

func BooleanDefaults() []bool {
	return []bool{false, true}
}

func BooleanUnique(idx int) bool {
	return idx%2 != 0
}

func NumericDefaults() []pgtype.Numeric {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(147454), nil)
	max.Add(max, big.NewInt(1))
	min := new(big.Int).Neg(max)

	return []pgtype.Numeric{
		{Int: big.NewInt(0), Status: pgtype.Present},
		{Int: big.NewInt(123), Exp: -2, Status: pgtype.Present},
		{Int: big.NewInt(-123), Exp: -2, Status: pgtype.Present},
		{Status: pgtype.Present, NaN: true},
		{Int: max, Exp: -16383, Status: pgtype.Present},
		{Int: min, Exp: -16383, Status: pgtype.Present},
	}

}

func NumericUnique(num int) pgtype.Numeric {
	n := pgtype.Numeric{}
	n.Set(fmt.Sprintf("%d.%d", num, num))
	return n
}

func TextDefaults() []string {
	return []string{"hello world", "3?!-+@.(\x01)ñ水불ツ😂"}
}

func TextUnique(num int) string {
	return fmt.Sprintf("unique_%d", num)
}

func TimestampTZDefaults() []pgtype.Timestamptz {
	// create normal timestamptz i.e. 1991-11-25 12:34:56.123456+07
	tz, _ := time.Parse(time.RFC3339Nano, "1991-11-25T06:34:56.123456+01:00")

	// create pg min range timestamptz i.e. 4714-11-24 00:22:00+00:22 BC
	// i.e. []byte{1, 255, 255, 255, 221, 94, 237, 229, 0, 0, 0, 0, 0, 0, 82}
	pgMin := new(time.Time)
	pgMin.UnmarshalBinary([]byte{1, 255, 255, 255, 221, 94, 237, 229, 0, 0, 0, 0, 0, 0, 82})

	// create pg max range timestamp i.e. 294276-12-31 23:59:59.999999+00
	// i.e. []byte{1, 0, 0, 8, 114, 43, 196, 208, 255, 59, 154, 198, 24, 0, 60}
	pgMax := new(time.Time)
	pgMax.UnmarshalBinary([]byte{1, 0, 0, 8, 114, 43, 196, 208, 255, 59, 154, 198, 24, 0, 60})

	return []pgtype.Timestamptz{
		{Time: tz.UTC(), Status: pgtype.Present},
		{Time: (*pgMin).UTC(), Status: pgtype.Present},
		{Time: (*pgMax).UTC(), Status: pgtype.Present},
		{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
}

func TimestamptTZUnique(num int) pgtype.Timestamptz {
	tz, _ := time.Parse(time.RFC3339Nano, "2000-01-02T01:23:45.123456+00:00")
	tz = tz.AddDate(num, 0, 0).UTC()
	return pgtype.Timestamptz{Time: tz, Status: pgtype.Present}
}

func DateDefaults() []pgtype.Date {
	t1, _ := time.Parse("2006-01-02", "1991-11-11")
	t2, _ := time.Parse("2006-01-02", "0001-11-24")
	t2 = t2.AddDate(-4714, 0, 0)
	t3, _ := time.Parse("2006-01-02", "0001-12-31")
	t3 = t3.AddDate(5874896, 0, 0)

	return []pgtype.Date{
		{Time: t1, Status: pgtype.Present},
		{Time: t2, Status: pgtype.Present},
		{Time: t3, Status: pgtype.Present},
		{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
}

func DateUnique(num int) pgtype.Date {
	t1, _ := time.Parse("2006-01-02", "2000-01-02")
	t1 = t1.AddDate(num, 0, 0)
	return pgtype.Date{Time: t1, Status: pgtype.Present}
}

func ByteDefaults() []pgtype.Bytea {
	return []pgtype.Bytea{
		{Status: pgtype.Present, Bytes: []byte("hello")},
		{Status: pgtype.Present, Bytes: []byte("mañana €5,90")},
		{Status: pgtype.Present, Bytes: []byte{0}},
	}
}
func ByteUnique(num int) pgtype.Bytea {
	return pgtype.Bytea{Status: pgtype.Present, Bytes: []byte(fmt.Sprintf("unique_%d", num))}
}

func UUID() []pgtype.UUID {
	return []pgtype.UUID{
		// "00010203-0405-0607-0809-0a0b0c0d0e0f"
		{
			Bytes:  [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			Status: pgtype.Present,
		},
		// "00000000-0000-0000-0000-000000000000"
		{
			Bytes:  [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Status: pgtype.Present,
		},
	}
}
