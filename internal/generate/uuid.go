package generate

import (
	"github.com/jackc/pgtype"
)

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
