package pgtestvals

import (
	"math/big"
	"time"

	"github.com/jackc/pgtype"
)

var (
	integerTestVals      = []interface{}{0, 1, -1, 2147483647, -2147483648}
	boolTestVals         = []interface{}{true, false}
	numericTestVals      = numeric()
	textTestVals         = []interface{}{"hello world", "3?!-+@.(\x01)ñ水불ツ😂"}
	timestatmptzTestVals = timestamptz()
	dateTestVals         = date()
	byteTestVals         = bytea()
	uuidTestVals         = uuid()
)

func numeric() []interface{} {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(147454), nil)
	max.Add(max, big.NewInt(1))
	min := new(big.Int).Neg(max)

	return []interface{}{
		pgtype.Numeric{Int: big.NewInt(0), Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(123), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(-123), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Status: pgtype.Present, NaN: true},
		pgtype.Numeric{Int: max, Exp: -16383, Status: pgtype.Present},
		pgtype.Numeric{Int: min, Exp: -16383, Status: pgtype.Present},
	}
}

func timestamptz() []interface{} {
	return []interface{}{
		pgtype.Timestamptz{Time: time.Date(1991, 11, 25, 5, 34, 56, 123456000, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Time: time.Date(294276, 12, 31, 23, 59, 59, 999999000, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
}

func date() []interface{} {
	return []interface{}{
		pgtype.Date{Time: time.Date(1991, 11, 11, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Time: time.Date(5874897, 12, 31, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
}

func bytea() []interface{} {
	return []interface{}{
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte("hello")},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte("mañana €5,90")},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{0}},
	}
}

func uuid() []interface{} {
	return []interface{}{
		// "00010203-0405-0607-0809-0a0b0c0d0e0f"
		pgtype.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		// "00000000-0000-0000-0000-000000000000"
		pgtype.UUID{Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Status: pgtype.Present},
	}
}
