// package generate provides fake values of specified types
package generate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/justasable/pgmock/internal/query"
)

// GenerateData generates mock data for all given tables in the db
func GenerateData(conn *pgx.Conn) error {
	tt, err := query.Tables(conn)
	if err != nil {
		return err
	}

	// first pass: generate rows for column types
	for _, t := range tt {
		err := generateForTableValues(conn, &t)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateForTableValues(conn *pgx.Conn, table *query.Table) error {
	// create task for each column
	var tt []*colTask
	for i := 0; i < len(table.Columns); i++ {
		t, err := newColTask(table.Columns[i])
		if err != nil {
			return err
		} else if t == nil {
			continue
		}
		tt = append(tt, t)
	}

	// generate data row until all tasks finished
	var errs []error
	for {
		var finished bool = true
		var colNames []string
		var values []interface{}

		// gather values for row
		for _, t := range tt {
			// check if we skip current value gen i.e. default val
			if t.ShouldSkip() {
				t.Advance()
				continue
			}

			colNames = append(colNames, t.column.Name)
			values = append(values, t.CurrentVal())

			if !t.Done() {
				finished = false
			}
			t.Advance()
		}

		if finished {
			break
		}

		// build db command
		var placeholders []string
		for idx := range values {
			placeholders = append(placeholders, fmt.Sprintf("$%d", idx+1))
		}
		cmd := fmt.Sprintf(
			"INSERT INTO %s.%s (%s) VALUES (%s)",
			table.Namespace, table.Name,
			strings.Join(colNames, ","),
			strings.Join(placeholders, ","),
		)

		// insert row
		_, err := conn.Exec(context.Background(), cmd, values...)
		if err != nil {
			// we don't return on this error as some unforseen constraints may prevent
			// row insert but we can still continue with the next row
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		var errStr string
		for _, err := range errs {
			errStr += fmt.Sprintf("%s\n", err)
		}
		return errors.New(errStr)
	}

	return nil
}
