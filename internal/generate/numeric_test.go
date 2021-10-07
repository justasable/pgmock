package generate_test

import (
	"math/big"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/stretchr/testify/assert"
)

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
