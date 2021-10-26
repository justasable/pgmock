package generate_test

import (
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/justasable/pgmock/internal/query"
	"github.com/stretchr/testify/assert"
)

func TestCompositeGenTestVals(t *testing.T) {
	// supported type (no default, not null) -> test vals
	col := query.Column{DataType: pgtype.Int4OID, IsNotNull: true}
	gen := generate.NewDataGenerator(col)
	assert.Equalf(t, []interface{}{0, 1, -1, 2147483647, -2147483648}, gen.TestVals(), "supported type (no default, not null)")

	// supported type (no default, nullable) -> null, test vals
	col = query.Column{DataType: pgtype.Int4OID}
	gen = generate.NewDataGenerator(col)
	assert.Equalf(t, []interface{}{nil, 0, 1, -1, 2147483647, -2147483648}, gen.TestVals(), "supported type (no default, nullable)")

	// supported type (has default, not null) -> test vals, then default
	col = query.Column{DataType: pgtype.Int4OID, HasDefault: true, IsNotNull: true}
	expected := []interface{}{0, 1, -1, 2147483647, -2147483648, generate.DefaultValType{}}
	gen = generate.NewDataGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "supported type (has default, not null)")

	// supported type (has default, nullable) -> null, test vals
	col = query.Column{DataType: pgtype.Int4OID, HasDefault: true}
	expected = []interface{}{nil, 0, 1, -1, 2147483647, -2147483648, generate.DefaultValType{}}
	gen = generate.NewDataGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "supported type (has default, nullable)")

	// Boolean (has default, nullable) -> null, test vals
	// special case as boolean test vals are exhaustive we skip default value that could case an insert error
	col = query.Column{DataType: pgtype.BoolOID, HasDefault: true}
	expected = []interface{}{nil, true, false}
	gen = generate.NewDataGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "bool (has default, nullable)")

}
