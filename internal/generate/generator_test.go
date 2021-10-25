package generate_test

import (
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/justasable/pgmock/internal/query"
	"github.com/stretchr/testify/assert"
)

func TestDefaultGenTestVals(t *testing.T) {
	c := query.Column{Generated: query.GENERATED_STORED}
	gen := generate.NewValueGenerator(c)
	got := gen.TestVals()
	expected := []interface{}{generate.DEFAULT_VAL}
	assert.Equal(t, expected, got)
}

func TestDefaultGenUniqueVal(t *testing.T) {
	c := query.Column{Generated: query.GENERATED_STORED}
	gen := generate.NewValueGenerator(c)
	got := gen.UniqueVal(0)
	assert.Equal(t, generate.DEFAULT_VAL, got)
}

func TestCompositeGenTestVals(t *testing.T) {
	// no default, not null
	col := query.Column{DataType: pgtype.BoolOID, HasDefault: false, IsNotNull: true}
	gen := generate.NewValueGenerator(col)
	assert.Equalf(t, []interface{}{true, false}, gen.TestVals(), "no default, not null")

	// default, not null
	col = query.Column{DataType: pgtype.BoolOID, HasDefault: true, IsNotNull: true}
	expected := []interface{}{true, false, generate.DEFAULT_VAL}
	gen = generate.NewValueGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "default, not null")

	// default, null
	col = query.Column{DataType: pgtype.BoolOID, HasDefault: true, IsNotNull: false}
	expected = []interface{}{nil, true, false, generate.DEFAULT_VAL}
	gen = generate.NewValueGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "default, nullable")
}
