package generate

import (
	"fmt"

	"github.com/justasable/pgmock/internal/query"
)

// colTask is a stateful data generation unit scoped to single column
type colTask struct {
	column    query.Column
	cursor    int
	generator Generator
}

// newColTask initialises a column task, returning nil if skipped
// and error if value generation for a column is not supported
func newColTask(c query.Column) (*colTask, error) {
	// check if we can skip column
	if c.Identity == query.IDENTITY_ALWAYS ||
		c.Constraint == query.CONSTRAINT_FOREIGN_KEY ||
		c.Generated == query.GENERATED_STORED {
		return nil, nil
	}

	// create column task
	gen, err := NewGenerator(c.DataType)
	if err != nil && !c.HasDefault && c.IsNotNull {
		return nil, fmt.Errorf("%w: column %s of type %d", err, c.Name, c.DataType)
	}

	t := &colTask{column: c, generator: gen}
	return t, nil
}

func (t *colTask) ShouldSkip() bool {
	return t.cursor == 0 && t.column.HasDefault
}

func (t *colTask) CurrentVal() interface{} {
	// return db default value if possible
	// return null value if possible
	// return our own generated values
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

	return t.generator.ValueForRow(idx)
}

func (t *colTask) Advance() {
	t.cursor++
}

func (t *colTask) Done() bool {
	return t.generator.Done(t.cursor)
}
