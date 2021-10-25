package generate

import (
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgtype"
)

// ValueGenerator generates test values and additional values unique to test value
type ValueGenerator interface {
	TestVals() []interface{}
	UniqueVal(int) interface{}
}

type integerGen struct{}

func (i integerGen) TestVals() []interface{} {
	return []interface{}{0, 1, -1, 2147483647, -2147483648}
}

func (i integerGen) UniqueVal(n int) interface{} {
	return 100 + n
}

type booleanGen struct{}

func (b booleanGen) TestVals() []interface{} {
	return []interface{}{true, false}
}

func (b booleanGen) UniqueVal(_ int) interface{} {
	return nil
}

type numericGen struct{}

func (n numericGen) TestVals() []interface{} {
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

func (n numericGen) UniqueVal(idx int) interface{} {
	num := pgtype.Numeric{}
	num.Set(fmt.Sprintf("%d.%d", idx, idx))
	return num
}

type textGen struct{}

func (t textGen) TestVals() []interface{} {
	return []interface{}{"hello world", "3?!-+@.(\x01)ñ水불ツ😂"}
}

func (t textGen) UniqueVal(n int) interface{} {
	return fmt.Sprintf("unique_%d", n)
}

type timestampTZGen struct{}

func (t timestampTZGen) TestVals() []interface{} {
	return []interface{}{
		pgtype.Timestamptz{Time: time.Date(1991, 11, 25, 5, 34, 56, 123456000, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Time: time.Date(294276, 12, 31, 23, 59, 59, 999999000, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
}

func (t timestampTZGen) UniqueVal(n int) interface{} {
	return pgtype.Timestamptz{Time: time.Date(2000+n, 1, 2, 1, 23, 45, 123456000, time.UTC), Status: pgtype.Present}
}

type dateGen struct{}

func (d dateGen) TestVals() []interface{} {
	return []interface{}{
		pgtype.Date{Time: time.Date(1991, 11, 11, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Time: time.Date(5874897, 12, 31, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
}

func (d dateGen) UniqueVal(n int) interface{} {
	return pgtype.Date{Time: time.Date(2000+n, 1, 2, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}
}

type byteGen struct{}

func (b byteGen) TestVals() []interface{} {
	return []interface{}{
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte("hello")},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte("mañana €5,90")},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{0}},
	}
}

func (b byteGen) UniqueVal(n int) interface{} {
	return pgtype.Bytea{Status: pgtype.Present, Bytes: []byte(fmt.Sprintf("unique_%d", n))}
}

type uuidGen struct{}

func (u uuidGen) TestVals() []interface{} {
	return []interface{}{
		// "00010203-0405-0607-0809-0a0b0c0d0e0f"
		pgtype.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: pgtype.Present},
		// "00000000-0000-0000-0000-000000000000"
		pgtype.UUID{Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Status: pgtype.Present},
	}
}

func (u uuidGen) UniqueVal(n int) interface{} {
	ret := new(pgtype.UUID)
	ret.Set(fmt.Sprintf("%0.32x", n))
	return *ret
}

type DefaultValType string

const DEFAULT_VAL DefaultValType = "DEFAULT"

type defaultGen struct{}

func (u defaultGen) TestVals() []interface{} {
	return []interface{}{DEFAULT_VAL}
}

func (u defaultGen) UniqueVal(_ int) interface{} {
	return DEFAULT_VAL
}

type compositeGen struct {
	prependVals []interface{}
	appendVals  []interface{}
	ValueGenerator
}

func (c compositeGen) TestVals() []interface{} {
	var ret []interface{}
	ret = append(ret, c.prependVals...)
	ret = append(ret, c.ValueGenerator.TestVals()...)
	ret = append(ret, c.appendVals...)
	return ret
}

func (c compositeGen) UniqueVal(n int) interface{} {
	return c.ValueGenerator.UniqueVal(n)
}
