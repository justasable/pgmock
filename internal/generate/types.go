package generate

import (
	"math/big"

	"github.com/jackc/pgtype"
)

func Integer() []int {
	return []int{0, 1, -1, 2147483648, -2147483648}
}

func Numeric() []pgtype.Numeric {
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

func Text() []string {
	return []string{"hello world", "3?!-+@.(\x01)ñ水불ツ😂"}
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
