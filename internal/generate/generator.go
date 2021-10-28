package generate

import (
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/pgtestvals"
	"github.com/justasable/pgmock/internal/query"
)

type defaultValType struct{}

// generator is a COLUMN SCOPED, STATEFUL data generator for postgresql context sensitive
// test values and unique values that takes into account DATA TYPE and COLUMN CONSTRAINTS
// the provided record is needed for obtaining foreign key references
type generator struct {
	column   query.Column
	record   *recordSet
	cursor   int
	testVals []interface{}
	uniqueFn func(int) interface{}
}

func newGenerator(c query.Column, r *recordSet) *generator {
	// Generated, Identity Always -> [default], then default, default...
	if c.Generated == query.GENERATED_STORED || c.Identity == query.IDENTITY_ALWAYS {
		return &generator{
			column:   c,
			record:   r,
			testVals: []interface{}{defaultValType{}},
			uniqueFn: func(int) interface{} { return defaultValType{} },
		}
	}

	// foreign key columns
	if c.Constraint == query.CONSTRAINT_FOREIGN_KEY {
		return newFKGenerator(c, r)
	}

	/*
	   Supported Types

	   Depending on column constraints we prepend/append values, the order is somewhat important
	   nil comes first, our generated test values second, then any database default value last

	   The reason is eg in a bool DEFAULT TRUE UNIQUE column, default val can clash with our test values
	   causing an error and preventing other test values from being inserted. Hence this order delays any
	   potential errors up to the inevitable moment, creating a greater chance for successful row generation
	*/

	// supported types
	var ret = &generator{column: c, record: r}
	ret.testVals = pgtestvals.TestVals(c.DataType)
	ret.uniqueFn = pgtestvals.UniqueFn(c.DataType)
	if ret.testVals != nil && ret.uniqueFn != nil {
		if !c.IsNotNull {
			ret.testVals = append([]interface{}{nil}, ret.testVals...)
		}
		if c.HasDefault && c.DataType != pgtype.BoolOID {
			// special case as boolean test vals are exhaustive
			// we skip default value that could case an insert error
			ret.testVals = append(ret.testVals, defaultValType{})
		}
		return ret
	}

	// unsupported type
	if c.HasDefault && c.IsNotNull {
		// --(has default, not null) -> [default], then default...
		ret.testVals = []interface{}{defaultValType{}}
		ret.uniqueFn = func(int) interface{} { return defaultValType{} }
	} else if c.HasDefault && !c.IsNotNull {
		// --(has default, nullable) -> [null, default], then default...
		ret.testVals = []interface{}{nil, defaultValType{}}
		ret.uniqueFn = func(int) interface{} { return defaultValType{} }
	} else if !c.IsNotNull {
		// --(no default, nullable) -> [null], then null, null...
		ret.testVals = []interface{}{nil}
		ret.uniqueFn = func(int) interface{} { return nil }
	} else {
		// --(no default, not null) -> cannot generate values
		return nil
	}

	return ret
}

func newFKGenerator(c query.Column, r *recordSet) *generator {
	var ret = &generator{column: c, record: r}
	pkey := r.PKeyForTable(c.FKTableID)

	// we skip column generation until we have valid fk reference
	if pkey == nil {
		return nil
	}

	// determine order of referenced column
	var conkeyPosition int = -1
	for idx, order := range c.ConKeys {
		if c.Order == order {
			conkeyPosition = idx
		}
	}
	if conkeyPosition == -1 || len(c.ConKeys) != len(c.FKColumns) {
		fmt.Printf("error finding column %s of order %d in conkeys %v with confkeys %v",
			c.Name, c.Order, c.ConKeys, c.FKColumns)
		return nil
	}

	// get value of pk column being referenced
	pkOrder := c.FKColumns[conkeyPosition]
	var refValue interface{}
	for _, record := range pkey {
		if record.Order == pkOrder {
			refValue = record.Value
		}
	}

	if !c.IsNotNull && refValue != nil {
		ret.testVals = []interface{}{nil, refValue}
	} else {
		ret.testVals = []interface{}{refValue}
	}
	ret.uniqueFn = func(int) interface{} { return refValue }

	return ret
}

func (g generator) currentVal() interface{} {
	if g.cursor < len(g.testVals) {
		return g.testVals[g.cursor]
	}
	return g.uniqueFn(g.cursor - len(g.testVals))
}

func (g *generator) advance() {
	g.cursor++
}

func (g generator) done() bool {
	return g.cursor >= len(g.testVals)
}
