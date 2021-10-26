// package generate provides fake values of specified types
package generate

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/justasable/pgmock/internal/query"
)

// GenerateData generates mock data for all given tables in the db
func GenerateData(conn *pgx.Conn) error {
	// fetch all tables
	tt, err := query.Tables(conn)
	if err != nil {
		return err
	}

	// generate data for each table
	for _, t := range tt {
		err := generateDataForTable(conn, t)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

// generateDataForTable generates rows of data until all column testvals are exhausted
// errors are not fatal, row data may still have been successfully generated
func generateDataForTable(conn *pgx.Conn, t query.Table) error {
	// generate column tasks
	var tasks []*colTask
	for _, col := range t.Columns {
		generator := NewDataGenerator(col)
		if generator == nil {
			return nil // can't generate data for this column, skip table
		}

		aTask := newColumnTask(col, generator)
		tasks = append(tasks, aTask)
	}

	// iterate through column tasks until done
	var errs []error
	for {
		// exit if all test vals exhausted
		var done bool = true
		for _, aTask := range tasks {
			if !aTask.done() {
				done = false
				break
			}
		}
		if done {
			break
		}

		// insert data
		builder := NewQueryBuilder(t.Namespace, t.Name)
		for _, task := range tasks {
			builder.AddColumnData(task.column.Name, task.currentVal())
			task.advance()
		}
		err := builder.RunQuery(conn)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors inserting data for table %s.%s: %+v", t.Namespace, t.Name, errs)
	}

	return nil
}
