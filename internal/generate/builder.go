package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
)

type builder struct {
	schemaName string
	tableName  string
	colNames   []string
	values     []interface{}
}

func newQueryBuilder(schemaName string, tableName string) *builder {
	return &builder{schemaName: schemaName, tableName: tableName}
}

func (b *builder) AddColumnData(name string, value interface{}) {
	b.colNames = append(b.colNames, name)
	b.values = append(b.values, value)
}

func (b builder) RunQuery(conn *pgx.Conn) ([]interface{}, error) {
	// strip away DEFAULT and NULL values from parameters list
	var count int
	var params []string
	var insertVals []interface{}
	for _, val := range b.values {
		if _, ok := val.(defaultValType); ok {
			params = append(params, "DEFAULT")
		} else if val == nil {
			params = append(params, "NULL")
		} else {
			count++
			params = append(params, fmt.Sprintf("$%d", count))
			insertVals = append(insertVals, val)
		}
	}

	query := fmt.Sprintf(
		"INSERT INTO %s.%s (%s) VALUES (%s) RETURNING %s",
		b.schemaName, b.tableName,
		strings.Join(b.colNames, ","),
		strings.Join(params, ","),
		strings.Join(b.colNames, ","),
	)

	// execute query and scan results
	totalColumns := len(b.colNames)
	values := make([]interface{}, totalColumns)
	valuePtrs := make([]interface{}, totalColumns)
	for i := 0; i < totalColumns; i++ {
		valuePtrs[i] = &values[i]
	}
	err := conn.
		QueryRow(context.Background(), query, insertVals...).
		Scan(valuePtrs...)

	return values, err
}
