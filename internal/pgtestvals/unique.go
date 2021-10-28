package pgtestvals

import (
	"fmt"
	"time"

	"github.com/jackc/pgtype"
)

func integerUniqueFn(n int) interface{} {
	return 100 + n
}

func boolUniqueFn(int) interface{} {
	return nil
}

func numericUniqueFn(n int) interface{} {
	num := pgtype.Numeric{}
	num.Set(fmt.Sprintf("%d.%d", n, n))
	return num
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
