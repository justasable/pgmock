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
	expected := []int{0, 1, -1, 2147483647, -2147483648}
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

func TestTimestampTZDefaults(t *testing.T) {
	expected := []pgtype.Timestamptz{
		{Time: time.Date(1991, 11, 25, 5, 34, 56, 123456000, time.UTC), Status: pgtype.Present},
		{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		{Time: time.Date(294276, 12, 31, 23, 59, 59, 999999000, time.UTC), Status: pgtype.Present},
		{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
	got := generate.TimestampTZDefaults()
	require.Len(t, got, len(expected))

	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], got[i])
	}
}

func TestTimestampTZUnique(t *testing.T) {
	got := generate.TimestamptTZUnique(100)
	assert.Equal(t, pgtype.Timestamptz{
		Time:   time.Date(2100, 1, 2, 1, 23, 45, 123456000, time.UTC),
		Status: pgtype.Present},
		got)
}

func TestDateDefaults(t *testing.T) {
	expected := []pgtype.Date{
		{Time: time.Date(1991, 11, 11, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		{Time: time.Date(-4713, 11, 24, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		{Time: time.Date(5874897, 12, 31, 0, 0, 0, 0, time.UTC), Status: pgtype.Present},
		{Status: pgtype.Present, InfinityModifier: pgtype.Infinity},
		{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity},
	}
	got := generate.DateDefaults()
	require.Len(t, got, len(expected))

	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], got[i])
	}
}

func TestDateUnique(t *testing.T) {
	expected := pgtype.Date{Time: time.Date(2100, 1, 2, 0, 0, 0, 0, time.UTC), Status: pgtype.Present}
	got := generate.DateUnique(100)
	assert.Equal(t, expected, got)
}

func TestByteDefaults(t *testing.T) {
	expected := []pgtype.Bytea{
		{Status: pgtype.Present, Bytes: []byte{104, 101, 108, 108, 111}},
		{Status: pgtype.Present, Bytes: []byte{109, 97, 195, 177, 97, 110, 97, 32, 226, 130, 172, 53, 44, 57, 48}},
		{Status: pgtype.Present, Bytes: []byte{0}},
	}
	got := generate.ByteDefaults()
	assert.Equal(t, expected, got)
}

func TestByteUnique(t *testing.T) {
	got := generate.ByteUnique(2345)
	assert.Equal(t, "unique_2345", string(got.Bytes))
	assert.Equal(t, pgtype.Present, got.Status)
}

func TestUUIDDefaults(t *testing.T) {
	expected := []string{
		"00010203-0405-0607-0809-0a0b0c0d0e0f",
		"00000000-0000-0000-0000-000000000000",
	}
	got := generate.UUIDDefaults()
	for idx, e := range expected {
		assert.Equal(t, strings.Replace(e, "-", "", -1), hex.EncodeToString(got[idx].Bytes[:]))
		assert.Equal(t, pgtype.Present, got[idx].Status)
	}
}

func TestUUIDUnique(t *testing.T) {
	got := generate.UUIDUnique(2345)
	assert.Equal(t, strings.Replace("00000000-0000-0000-0000-000000000929", "-", "", -1), hex.EncodeToString(got.Bytes[:]))
	assert.Equal(t, pgtype.Present, got.Status)
}
