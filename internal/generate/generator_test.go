package generate_test

import (
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/justasable/pgmock/internal/query"
	"github.com/stretchr/testify/assert"
)

func TestUnsupportedTypeTestVals(t *testing.T) {
	// generated column
	col := query.Column{Generated: query.GENERATED_STORED}
	gen := generate.Test_newGenerator(col, nil)
	assert.Equal(t, []interface{}{generate.Test_defaultValType{}}, gen.Test_TestVals())
	assert.Equal(t, generate.Test_defaultValType{}, gen.Test_UniqueVal(0))

	// unsupported type (no default, not null) -> cannot generate
	col = query.Column{IsNotNull: true}
	gen = generate.Test_newGenerator(col, nil)
	assert.Nil(t, gen)

	// unsupported type (no default, nullable) --> null, then null, null...
	col = query.Column{}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equal(t, []interface{}{nil}, gen.Test_TestVals())
	assert.Equal(t, nil, gen.Test_UniqueVal(1))

	// unsupported type (has default, not null) -> default, then default, default...
	col = query.Column{HasDefault: true, IsNotNull: true}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equal(t, []interface{}{generate.Test_defaultValType{}}, gen.Test_TestVals())
	assert.Equal(t, generate.Test_defaultValType{}, gen.Test_UniqueVal(0))

	// unsupported type (has default, nullable) -> nil, default, then default, default...
	col = query.Column{HasDefault: true, IsNotNull: false}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equal(t, []interface{}{nil, generate.Test_defaultValType{}}, gen.Test_TestVals())
	assert.Equal(t, generate.Test_defaultValType{}, gen.Test_UniqueVal(0))
}

func TestConstraintsTestVals(t *testing.T) {
	// supported type (no default, not null) -> test vals
	col := query.Column{DataType: pgtype.Int4OID, IsNotNull: true}
	gen := generate.Test_newGenerator(col, nil)
	assert.Equalf(t, []interface{}{0, 1, -1, 2147483647, -2147483648}, gen.Test_TestVals(), "supported type (no default, not null)")

	// supported type (no default, nullable) -> null, test vals
	col = query.Column{DataType: pgtype.Int4OID}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equalf(t, []interface{}{nil, 0, 1, -1, 2147483647, -2147483648}, gen.Test_TestVals(), "supported type (no default, nullable)")

	// supported type (has default, not null) -> test vals, then default
	col = query.Column{DataType: pgtype.Int4OID, HasDefault: true, IsNotNull: true}
	expected := []interface{}{0, 1, -1, 2147483647, -2147483648, generate.Test_defaultValType{}}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equalf(t, expected, gen.Test_TestVals(), "supported type (has default, not null)")

	// supported type (has default, nullable) -> null, test vals
	col = query.Column{DataType: pgtype.Int4OID, HasDefault: true}
	expected = []interface{}{nil, 0, 1, -1, 2147483647, -2147483648, generate.Test_defaultValType{}}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equalf(t, expected, gen.Test_TestVals(), "supported type (has default, nullable)")

	// Boolean (has default, nullable) -> null, test vals
	// special case as boolean test vals are exhaustive we skip default value that could case an insert error
	col = query.Column{DataType: pgtype.BoolOID, HasDefault: true}
	expected = []interface{}{nil, true, false}
	gen = generate.Test_newGenerator(col, nil)
	assert.Equalf(t, expected, gen.Test_TestVals(), "bool (has default, nullable)")
}
