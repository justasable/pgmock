package generate_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerDefaults(t *testing.T) {
	expected := []int{0, 1, -1, 2147483648, -2147483648}
	assert.Equal(t, expected, generate.IntegerDefaults())
}

func TestIntegerUnique(t *testing.T) {
	assert.Equal(t, 105, generate.IntegerUnique(5))
}

func TestBooleanDefaults(t *testing.T) {
	assert.Equal(t, []bool{false, true}, generate.BooleanDefaults())
}

func TestBooleanUnique(t *testing.T) {
	assert.Equal(t, false, generate.BooleanUnique(50))
	assert.Equal(t, true, generate.BooleanUnique(51))
}

func TestNumericDefaults(t *testing.T) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(147454), nil)
	max.Add(max, big.NewInt(1))
	min := new(big.Int).Neg(max)

	expected := []pgtype.Numeric{
		{Int: big.NewInt(0), Status: pgtype.Present},
		{Int: big.NewInt(123), Exp: -2, Status: pgtype.Present},
		{Int: big.NewInt(-123), Exp: -2, Status: pgtype.Present},
		{Status: pgtype.Present, NaN: true},
		{Int: max, Exp: -16383, Status: pgtype.Present},
		{Int: min, Exp: -16383, Status: pgtype.Present},
	}

	for idx, e := range expected {
		assert.Equalf(t, e, generate.NumericDefaults()[idx], "mismatched value at index %d", idx)
	}
}

func TestNumericUnique(t *testing.T) {
	assert.Equal(t, pgtype.Numeric{Int: big.NewInt(7070), Exp: -2, Status: pgtype.Present}, generate.NumericUnique(70))
}

func TestTextDefaults(t *testing.T) {
	expected := []string{
		"hello world",
		"3?!-+@.(\x01)\u00f1\u6c34\ubd88\u30c4\U0001f602",
	}
	assert.Equal(t, expected, generate.TextDefaults())
}

func TestTextUnique(t *testing.T) {
	assert.Equal(t, "unique_1010", generate.TextUnique(1010))
}

func TestTimestampTZ(t *testing.T) {
	expected := []string{
		"2021-11-01T05:34:56.123456Z",
		"-4713-11-24T00:00:00Z",
		"294276-12-31T23:59:59.999999Z",
	}
	got := generate.TimestampTZ()

	require.Len(t, got, len(expected)+2)

	// check non infinity times
	for i := 0; i < 3; i++ {
		assert.Equal(t, expected[i], got[i].Time.Format(time.RFC3339Nano))
	}

	// check infinity times
	assert.Equal(t, pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.Infinity}, got[3])
	assert.Equal(t, pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity}, got[4])
}

func TestDate(t *testing.T) {
	expected := []string{
		"1991-11-11",
		"-4713-11-24",
		"5874897-12-31",
	}
	got := generate.Date()

	// check non infinity times
	for i := 0; i < 3; i++ {
		assert.Equal(t, expected[i], got[i].Time.Format("2006-01-02"))
		assert.Equal(t, pgtype.Present, got[i].Status)
		assert.Equal(t, pgtype.None, got[i].InfinityModifier)
	}

	// check infinity times
	assert.Equal(t, pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.Infinity}, got[3])
	assert.Equal(t, pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity}, got[4])
}

func TestByte(t *testing.T) {
	expected := []pgtype.Bytea{
		{Status: pgtype.Present, Bytes: []byte{104, 101, 108, 108, 111}},
		{Status: pgtype.Present, Bytes: []byte{109, 97, 195, 177, 97, 110, 97, 32, 226, 130, 172, 53, 44, 57, 48}},
		{Status: pgtype.Present, Bytes: []byte{0}},
	}
	got := generate.Byte()

	assert.Equal(t, expected, got)
}

func TestUUID(t *testing.T) {
	got := generate.UUID()
	expected := []string{
		"00010203-0405-0607-0809-0a0b0c0d0e0f",
		"00000000-0000-0000-0000-000000000000",
	}

	require.Len(t, got, len(expected))

	for k, v := range got {
		assert.Equalf(t, pgtype.Present, v.Status, "element at index %d status not set to present", k)
		assert.Equalf(t, strings.Replace(expected[k], "-", "", -1), hex.EncodeToString(v.Bytes[:]),
			"element at index %d values do not match", k)
	}
}
