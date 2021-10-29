package pgtestvals

import (
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgtype"
)

func integerUniqueFn(n int) interface{} {
	return 100 + n
}

func boolUniqueFn(int) interface{} {
	return nil
}

func numericUniqueFn(typeMod int) func(int) interface{} {
	// numeric
	if typeMod == -1 {
		return func(n int) interface{} {
			num := big.NewInt(10)
			num.Exp(num, big.NewInt(int64(n+n+1)), nil)
			num.Add(num, big.NewInt(1))
			return pgtype.Numeric{Int: num, Exp: -int32(n + 1), Status: pgtype.Present}
		}
	}

	// numeric(p, s)
	scale := (typeMod - 4) & 0xffff
	return func(n int) interface{} {
		return pgtype.Numeric{Int: big.NewInt(int64(n)), Exp: -int32(scale), Status: pgtype.Present}
	}
}

func textUniqueFn(n int) interface{} {
	return fmt.Sprintf("unique_%d", n)
}

func timestamptzUniqueFn(n int) interface{} {
	return pgtype.Timestamptz{Time: time.Date(2000+n, 1, 2, 1, 23, 45, 123456000, time.UTC), Status: pgtype.Present}
}

func dateUniqueFn(n int) interface{} {
	return pgtype.Date{Time: time.Date(2000+n, 1, 2, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}
}

func byteUniqueFn(n int) interface{} {
	return pgtype.Bytea{Status: pgtype.Present, Bytes: []byte(fmt.Sprintf("unique_%d", n))}
}

func uuidUniqueFn(n int) interface{} {
	ret := new(pgtype.UUID)
	ret.Set(fmt.Sprintf("%0.32x", n+1))
	return *ret
}
