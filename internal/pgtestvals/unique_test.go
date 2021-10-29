package pgtestvals_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/pgtestvals"
	"github.com/stretchr/testify/assert"
)

// convenience function for testing simple data type unique val generation
func testGenUniqueVal(t *testing.T, name string, dataType pgtype.OID, count int, expected interface{}) {
	fn := pgtestvals.UniqueFn(dataType, -1)
	got := fn(count)
	assert.Equalf(t, expected, got, "%s unique val gen mismatch for count %d", name, count)
}
func TestIntegerGenUniqueVal(t *testing.T) {
	testGenUniqueVal(t, "integer", pgtype.Int4OID, 5, 105)
}
func TestBooleanGenUniqueVal(t *testing.T) {
	testGenUniqueVal(t, "bool", pgtype.BoolOID, 1, nil)
}

func TestNumericGenUniqueVal(t *testing.T) {
	expected := pgtype.Numeric{Int: big.NewInt(100001), Exp: -3, Status: pgtype.Present}
	testGenUniqueVal(t, "numeric", pgtype.NumericOID, 2, expected)
}

func TestNumericPrecisionScaleGenUniqueVal(t *testing.T) {
	// numeric(5, 2) has attypmod of 327686
	expected := pgtype.Numeric{Int: big.NewInt(99999), Exp: -2, Status: pgtype.Present}
	fn := pgtestvals.UniqueFn(pgtype.NumericOID, 327686)
	got := fn(99999)
	assert.Equal(t, expected, got)
}

func TestTextGenUniqueVal(t *testing.T) {
	testGenUniqueVal(t, "text", pgtype.TextOID, 5, "unique_5")
}

func TestTimestampTZGenUniqueVal(t *testing.T) {
	expected := pgtype.Timestamptz{Time: time.Date(2100, 1, 2, 1, 23, 45, 123456000, time.UTC), Status: pgtype.Present}
	testGenUniqueVal(t, "timestamptz", pgtype.TimestamptzOID, 100, expected)
}

func TestDateGenUniqueVal(t *testing.T) {
	expected := pgtype.Date{Time: time.Date(2100, 1, 2, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}
	testGenUniqueVal(t, "date", pgtype.DateOID, 100, expected)
}

func TestByteGenUniqueVal(t *testing.T) {
	expected := pgtype.Bytea{Bytes: []byte("unique_2345"), Status: pgtype.Present}
	testGenUniqueVal(t, "byte", pgtype.ByteaOID, 2345, expected)
}

func TestUUIDGenUniqueVal(t *testing.T) {
	expected := new(pgtype.UUID)
	expected.Set("00000000-0000-0000-0000-000000000929")
	testGenUniqueVal(t, "uuid", pgtype.UUIDOID, 2344, *expected)

	expected.Set("00000000-0000-0000-0000-000000000001")
	testGenUniqueVal(t, "uuid", pgtype.UUIDOID, 0, *expected)
}
