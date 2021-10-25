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
	col := query.Column{DataType: pgtype.BoolOID, IsNotNull: true}
	gen := generate.NewValueGenerator(col)
	assert.Equalf(t, []interface{}{true, false}, gen.TestVals(), "supported type (no default, not null)")

	// supported type (no default, nullable) -> null, test vals
	col = query.Column{DataType: pgtype.BoolOID}
	gen = generate.NewValueGenerator(col)
	assert.Equalf(t, []interface{}{nil, true, false}, gen.TestVals(), "supported type (no default, nullable)")

	// supported type (has default, not null) -> test vals, then default
	col = query.Column{DataType: pgtype.BoolOID, HasDefault: true, IsNotNull: true}
	expected := []interface{}{true, false, generate.DEFAULT_VAL}
	gen = generate.NewValueGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "supported type (has default, not null)")

	// supported type (has default, nullable) -> null, test vals
	col = query.Column{DataType: pgtype.BoolOID, HasDefault: true}
	expected = []interface{}{nil, true, false, generate.DEFAULT_VAL}
	gen = generate.NewValueGenerator(col)
	assert.Equalf(t, expected, gen.TestVals(), "supported type (has default, nullable)")
}
