package generate_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/justasable/pgmock/internal/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// convenience function for testing simple data type test vals generation
func testGenTestVals(t *testing.T, name string, dataType pgtype.OID, expected []interface{}) {
	col := query.Column{DataType: dataType, IsNotNull: true}
	gen := generate.NewValueGenerator(col)
	got := gen.TestVals()

	require.Lenf(t, got, len(expected), "length mismatch for %s", name)
	for i := 0; i < len(expected); i++ {
		assert.Equalf(t, expected[i], got[i], "%s element at index %d does not match", name, i)
	}
}

// convenience function for testing simple data type unique val generation
func testGenUniqueVal(t *testing.T, name string, dataType pgtype.OID, count int, expected interface{}) {
	col := query.Column{DataType: dataType, IsNotNull: true}
	gen := generate.NewValueGenerator(col)
	got := gen.UniqueVal(count)
	assert.Equalf(t, expected, got, "%s unique val gen mismatch for count %d", name, count)
}

func TestIntegerGenTestVals(t *testing.T) {
	testGenTestVals(t, "integer", pgtype.Int4OID, []interface{}{0, 1, -1, 2147483647, -2147483648})
}

func TestIntegerGenUniqueVal(t *testing.T) {
	testGenUniqueVal(t, "integer", pgtype.Int4OID, 5, 105)
}

func TestBooleanGenTestVals(t *testing.T) {
	testGenTestVals(t, "bool", pgtype.BoolOID, []interface{}{true, false})
}

func TestBooleanGenUniqueVal(t *testing.T) {
	testGenUniqueVal(t, "bool", pgtype.BoolOID, 1, nil)
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
		pgtype.Numeric{Int: max, Exp: -16383, Status: pgtype.Present},
		pgtype.Numeric{Int: min, Exp: -16383, Status: pgtype.Present},
	}
	testGenTestVals(t, "numeric", pgtype.NumericOID, expected)
}

func TestNumericGenUniqueVal(t *testing.T) {
	expected := pgtype.Numeric{Int: big.NewInt(7070), Exp: -2, Status: pgtype.Present}
	testGenUniqueVal(t, "numeric", pgtype.NumericOID, 70, expected)
}

func TestTextGenTestVals(t *testing.T) {
	expected := []interface{}{
		"hello world",
		"3?!-+@.(\x01)\u00f1\u6c34\ubd88\u30c4\U0001f602",
	}
	testGenTestVals(t, "text", pgtype.TextOID, expected)
}

func TestTextGenUniqueVal(t *testing.T) {
	testGenUniqueVal(t, "text", pgtype.TextOID, 5, "unique_5")
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

func TestTimestampTZGenUniqueVal(t *testing.T) {
	expected := pgtype.Timestamptz{Time: time.Date(2100, 1, 2, 1, 23, 45, 123456000, time.UTC), Status: pgtype.Present}
	testGenUniqueVal(t, "timestamptz", pgtype.TimestamptzOID, 100, expected)
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

func TestDateGenUniqueVal(t *testing.T) {
	expected := pgtype.Date{Time: time.Date(2100, 1, 2, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}
	testGenUniqueVal(t, "date", pgtype.DateOID, 100, expected)
}

func TestByteGenTestVals(t *testing.T) {
	expected := []interface{}{
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{104, 101, 108, 108, 111}},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{109, 97, 195, 177, 97, 110, 97, 32, 226, 130, 172, 53, 44, 57, 48}},
		pgtype.Bytea{Status: pgtype.Present, Bytes: []byte{0}},
	}
	testGenTestVals(t, "byte", pgtype.ByteaOID, expected)
}

func TestByteGenUniqueVal(t *testing.T) {
	expected := pgtype.Bytea{Bytes: []byte("unique_2345"), Status: pgtype.Present}
	testGenUniqueVal(t, "byte", pgtype.ByteaOID, 2345, expected)
}

func TestUUIDGenTestVals(t *testing.T) {
	n := new(pgtype.UUID)
	n.Set("00010203-0405-0607-0809-0a0b0c0d0e0f")
	m := new(pgtype.UUID)
	m.Set("00000000-0000-0000-0000-000000000000")
	expected := []interface{}{*n, *m}
	testGenTestVals(t, "uuid", pgtype.UUIDOID, expected)
}

func TestUUIDGenUniqueVal(t *testing.T) {
	expected := new(pgtype.UUID)
	expected.Set("00000000-0000-0000-0000-000000000929")
	testGenUniqueVal(t, "uuid", pgtype.UUIDOID, 2345, *expected)
}

func TestUnsupportedType(t *testing.T) {
	// generated column
	col := query.Column{Generated: query.GENERATED_STORED}
	gen := generate.NewValueGenerator(col)
	assert.Equal(t, []interface{}{generate.DEFAULT_VAL}, gen.TestVals())
	assert.Equal(t, generate.DEFAULT_VAL, gen.UniqueVal(0))

	// unsupported type (no default, not null) -> cannot generate
	col = query.Column{IsNotNull: true}
	gen = generate.NewValueGenerator(col)
	assert.Nil(t, gen)

	// unsupported type (no default, nullable) --> null, then null, null...
	col = query.Column{}
	gen = generate.NewValueGenerator(col)
	assert.Equal(t, []interface{}{nil}, gen.TestVals())
	assert.Equal(t, nil, gen.UniqueVal(1))

	// unsupported type (has default, not null) -> default, then default, default...
	col = query.Column{HasDefault: true, IsNotNull: true}
	gen = generate.NewValueGenerator(col)
	assert.Equal(t, []interface{}{generate.DEFAULT_VAL}, gen.TestVals())
	assert.Equal(t, generate.DEFAULT_VAL, gen.UniqueVal(0))

	// unsupported type (has default, nullable) -> nil, default, then default, default...
	col = query.Column{HasDefault: true, IsNotNull: false}
	gen = generate.NewValueGenerator(col)
	assert.Equal(t, []interface{}{nil, generate.DEFAULT_VAL}, gen.TestVals())
	assert.Equal(t, generate.DEFAULT_VAL, gen.UniqueVal(0))

}
