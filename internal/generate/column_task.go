package generate

// import (
// 	"fmt"

// 	"github.com/justasable/pgmock/internal/query"
// )

// // colTask is a stateful data generation unit scoped to single column
// type colTask struct {
// 	column    query.Column
// 	cursor    int
// 	generator Generator
// }

// // newColTask initialises a column task
// // whose purpose is simple data type generation
// // returns nil data type not supported
// // and error if value generation for a column is not supported
// func newColTask(c query.Column) (*colTask, error) {
// 	// check if we can skip column
// 	if c.Identity == query.IDENTITY_ALWAYS ||
// 		c.Constraint == query.CONSTRAINT_FOREIGN_KEY ||
// 		c.Generated == query.GENERATED_STORED {
// 		return nil, nil
// 	}

// 	// create column task
// 	gen, err := NewGenerator(c.DataType)
// 	if err != nil && !c.HasDefault && c.IsNotNull {
// 		return nil, fmt.Errorf("%w: column %s of type %d", err, c.Name, c.DataType)
// 	}

// 	t := &colTask{column: c, generator: gen}
// 	return t, nil
// }

// func (t *colTask) CurrentVal() interface{} {
// 	// return db default value if possible
// 	if t.cursor == 0 && t.column.HasDefault {
// 		return DEFAULT_VAL
// 	}

// 	// return null value if possible
// 	if !t.column.IsNotNull && (t.cursor == 0 && !t.column.HasDefault || t.cursor == 1 && t.column.HasDefault) {
// 		return nil
// 	}

// 	// return our own generated values
// 	idx := t.cursor
// 	if t.column.HasDefault {
// 		idx--
// 	}
// 	if !t.column.IsNotNull {
// 		idx--
// 	}

// 	return t.generator.ValueForRow(idx)
// }

// func (t *colTask) Advance() {
// 	t.cursor++
// }

// func (t *colTask) Done() bool {
// 	return t.generator.Done(t.cursor)
// }
