package generate

import (
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/query"
)

// task is a data generation unit scoped to single column
type task struct {
	column *query.Column
	cursor int
}

// NewTask initialises a task, returning nil if skipped
// and error if value generation for a column is not supported
func NewTask(c *query.Column) (*task, error) {
	// check if we can skip column
	if c.Identity == query.IDENTITY_ALWAYS ||
		c.Constraint == query.CONSTRAINT_PRIMARY_KEY ||
		c.Constraint == query.CONSTRAINT_FOREIGN_KEY ||
		c.Generated == query.GENERATED_STORED {
		return nil, nil
	}

	// check if we support data type
	if !c.HasDefault && c.IsNotNull &&
		!(c.DataType == pgtype.Int4OID ||
			c.DataType == pgtype.BoolOID ||
			c.DataType == pgtype.NumericOID ||
			c.DataType == pgtype.TextOID ||
			c.DataType == pgtype.TimestamptzOID ||
			c.DataType == pgtype.DateOID ||
			c.DataType == pgtype.ByteaOID ||
			c.DataType == pgtype.UUIDOID) {
		// return error if we are unable to generate values for column
		return nil, fmt.Errorf("data type for column %s not supported", c.Name)
	}

	t := &task{column: c}
	return t, nil
}

func (t *task) ShouldSkip() bool {
	return t.cursor == 0 && t.column.HasDefault
}

func (t *task) CurrentVal() interface{} {
	if !t.column.IsNotNull && (t.cursor == 0 && !t.column.HasDefault || t.cursor == 1 && t.column.HasDefault) {
		return nil
	}

	idx := t.cursor
	if t.column.HasDefault {
		idx--
	}
	if !t.column.IsNotNull {
		idx--
	}

	switch t.column.DataType {
	case pgtype.Int4OID:
		vals := IntegerDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return IntegerUnique(idx)
	case pgtype.BoolOID:
		vals := BooleanDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return BooleanUnique(idx)
	case pgtype.NumericOID:
		vals := NumericDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return NumericUnique(idx)
	case pgtype.TextOID:
		vals := TextDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return TextUnique(idx)
	case pgtype.TimestamptzOID:
		vals := TimestampTZDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return TimestamptTZUnique(idx)
	case pgtype.DateOID:
		vals := DateDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return DateUnique(idx)
	case pgtype.ByteaOID:
		vals := ByteDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return ByteUnique(idx)
	case pgtype.UUIDOID:
		vals := UUIDDefaults()
		if idx < len(vals) {
			return vals[idx]
		}
		return UUIDUnique(idx)
	}

	return nil
}

func (t *task) Advance() {
	t.cursor++
}

func (t *task) Finished() bool {
	switch t.column.DataType {
	case pgtype.Int4OID:
		return t.cursor >= len(IntegerDefaults())
	case pgtype.BoolOID:
		return t.cursor >= len(BooleanDefaults())
	case pgtype.NumericOID:
		return t.cursor >= len(NumericDefaults())
	case pgtype.TextOID:
		return t.cursor >= len(TextDefaults())
	case pgtype.TimestamptzOID:
		return t.cursor >= len(TimestampTZDefaults())
	case pgtype.DateOID:
		return t.cursor >= len(DateDefaults())
	case pgtype.ByteaOID:
		return t.cursor >= len(ByteDefaults())
	case pgtype.UUIDOID:
		return t.cursor >= len(UUIDDefaults())
	}

	return true
}
