package pgtestvals_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/pgtestvals"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// convenience function for testing simple data type test vals generation
func testGenTestVals(t *testing.T, name string, dataType pgtype.OID, expected []interface{}) {
	got := pgtestvals.TestVals(dataType, -1)
	require.Lenf(t, got, len(expected), "length mismatch for %s", name)
	for i := 0; i < len(expected); i++ {
		assert.Equalf(t, expected[i], got[i], "%s element at index %d does not match", name, i)
	}
}

func TestIntegerGenTestVals(t *testing.T) {
	testGenTestVals(t, "integer", pgtype.Int4OID, []interface{}{0, 1, -1, 2147483647, -2147483648})
}

func TestBooleanGenTestVals(t *testing.T) {
	testGenTestVals(t, "bool", pgtype.BoolOID, []interface{}{true, false})
}

func TestNumericGenTestVals(t *testing.T) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(147454), nil)
	max.Add(max, big.NewInt(1))
	min := new(big.Int).Neg(max)
	expected := []interface{}{
		pgtype.Numeric{Int: big.NewInt(0), Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(123), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(-123), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Status: pgtype.Present, NaN: true},
		pgtype.Numeric{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Numeric{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
		pgtype.Numeric{Int: max, Exp: -16383, Status: pgtype.Present},
		pgtype.Numeric{Int: min, Exp: -16383, Status: pgtype.Present},
	}
	testGenTestVals(t, "numeric", pgtype.NumericOID, expected)
}

func TestNumericPrecisionScaleGenTestVals(t *testing.T) {
	// numeric(5, 2) has attypmod of 327686
	expected := []interface{}{
		pgtype.Numeric{Int: big.NewInt(0), Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(123), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(-123), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Status: pgtype.Present, NaN: true},
		pgtype.Numeric{Int: big.NewInt(99999), Exp: -2, Status: pgtype.Present},
		pgtype.Numeric{Int: big.NewInt(-99999), Exp: -2, Status: pgtype.Present},
	}

	got := pgtestvals.TestVals(pgtype.NumericOID, 327686)
	assert.Equal(t, expected, got)
}

func TestTextGenTestVals(t *testing.T) {
	expected := []interface{}{
		"hello world",
		"3?!-+@.(\x01)\u00f1\u6c34\ubd88\u30c4\U0001f602",
	}
	testGenTestVals(t, "text", pgtype.TextOID, expected)
}

func TestTimestampTZGenTestVals(t *testing.T) {
	expected := []interface{}{
		pgtype.Timestamptz{Time: time.Date(1991, 11, 25, 5, 34, 56, 123456000, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Time: time.Date(294276, 12, 31, 23, 59, 59, 999999000, time.UTC), Status: pgtype.Present},
		pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
	testGenTestVals(t, "timestamptz", pgtype.TimestamptzOID, expected)
}

func TestDateGenTestVals(t *testing.T) {
	expected := []interface{}{
		pgtype.Date{Time: time.Date(1991, 11, 11, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Time: time.Date(5874897, 12, 31, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
	testGenTestVals(t, "date", pgtype.DateOID, expected)
}

func TestByteGenTestVals(t *testing.T) {
	expected := []interface{}{
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{104, 101, 108, 108, 111}},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{109, 97, 195, 177, 97, 110, 97, 32, 226, 130, 172, 53, 44, 57, 48}},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{0}},
	}
	testGenTestVals(t, "byte", pgtype.ByteaOID, expected)
}

func TestUUIDGenTestVals(t *testing.T) {
	n := new(pgtype.UUID)
	n.Set("00010203-0405-0607-0809-0a0b0c0d0e0f")
	m := new(pgtype.UUID)
	m.Set("00000000-0000-0000-0000-000000000000")
	expected := []interface{}{*n, *m}
	testGenTestVals(t, "uuid", pgtype.UUIDOID, expected)
}
