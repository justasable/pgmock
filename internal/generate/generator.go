package generate

import (
	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/query"
)

// ValueGenerator generates test values and additional values unique to test value
type ValueGenerator interface {
	TestVals() []interface{}
	UniqueVal(int) interface{}
}

func NewValueGenerator(c query.Column) ValueGenerator {
	// generated columns
	if c.Generated == query.GENERATED_STORED {
		return defaultGen{}
	}

	// data type
	var ret compositeGen
	switch c.DataType {
	case pgtype.Int4OID:
		ret.ValueGenerator = integerGen{}
	case pgtype.BoolOID:
		ret.ValueGenerator = booleanGen{}
	case pgtype.NumericOID:
		ret.ValueGenerator = numericGen{}
	case pgtype.TextOID:
		ret.ValueGenerator = textGen{}
	case pgtype.TimestamptzOID:
		ret.ValueGenerator = timestampTZGen{}
	case pgtype.DateOID:
		ret.ValueGenerator = dateGen{}
	case pgtype.ByteaOID:
		ret.ValueGenerator = byteGen{}
	case pgtype.UUIDOID:
		ret.ValueGenerator = uuidGen{}
	default:
		// unsupported type (no default, not null) -> we're unable to generate anything
		if !c.HasDefault && c.IsNotNull {
			return nil
		}

		// unsupported type (no default, nullable) -> null, null, null...
		if !c.HasDefault && !c.IsNotNull {
			return nullGen{}
		}

		// unsupported type (has default, not null) -> default, default...
		if c.IsNotNull {
			return defaultGen{}
		}

		// unsupported type (has default, nullable) -> null, default, default...
		ret.prependVals = []interface{}{nil}
		ret.ValueGenerator = defaultGen{}
		return ret
	}

	/*
		Supported Types

		Depending on column constraints we prepend/append values, the order is somewhat important
		nil comes first, our generated test values second, then any database default value last

		The reason is eg in a bool DEFAULT TRUE UNIQUE column, default val can clash with our test values
		causing an error and preventing other test values from being inserted. Hence this order delays any
		potential errors up to the inevitable moment, creating a greater chance for successful row generation
	*/

	if !c.IsNotNull {
		ret.prependVals = []interface{}{nil}
	}
	if c.HasDefault {
		ret.appendVals = []interface{}{DEFAULT_VAL}
	}

	return ret
}
