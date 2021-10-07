package generate_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInteger(t *testing.T) {
	expected := []int{0, 1, -1, 2147483648, -2147483648}
	got := generate.Integer()

	assert.Equal(t, got, expected)
}

func TestNumeric(t *testing.T) {
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

	assert.Equal(t, expected, generate.Numeric())
}

func TestText(t *testing.T) {
	expected := []string{
		"hello world",
		"3?!-+@.(\x01)\u00f1\u6c34\ubd88\u30c4\U0001f602",
	}
	got := generate.Text()
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
