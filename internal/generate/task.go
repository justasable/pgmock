package generate

import (
	"github.com/justasable/pgmock/internal/query"
)

// colTask represents a stateful unit of work scoped to a column
type colTask struct {
	column   query.Column
	cursor   int
	testVals []interface{}
	uniqueFn func(int) interface{}
}

func newColumnTask(col query.Column, generator DataGenerator) *colTask {
	return &colTask{column: col, testVals: generator.TestVals(), uniqueFn: generator.UniqueVal}
}

func (c *colTask) currentVal() interface{} {
	// return testval
	if c.cursor < len(c.testVals) {
		return c.testVals[c.cursor]
	}

	// return unique val once testvals exhausted
	return c.uniqueFn(c.cursor - len(c.testVals))
}

func (c *colTask) advance() {
	c.cursor++
}

func (c *colTask) done() bool {
	return c.cursor >= len(c.testVals)
}
