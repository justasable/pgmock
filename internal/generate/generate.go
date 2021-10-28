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

	// keep attempting to generate data
	var record = &recordSet{}
	var count int
	for {
		// attempt for each table
		for _, t := range tt {
			if !record.HasData(t.ID) {
				generateDataForTable(conn, record, t)
			}
		}

		// exit condition
		totalRecords := record.TotalRecords()
		if count == totalRecords {
			break
		}
		count = totalRecords
	}

	return nil
}

// generateDataForTable attempts to generate rows of data until all column testvals are exhausted
// errors are not propagated as various db constraints can cause them which are outside
// our control, and we can still continue with data generation
func generateDataForTable(conn *pgx.Conn, r *recordSet, t query.Table) {
	// create column generators
	var colGens []*generator
	for _, col := range t.Columns {
		gen := newGenerator(col, r)
		// skip table if we're unable to generate data for a column
		if gen == nil {
			return
		}
		colGens = append(colGens, gen)
	}

	// iterate through column generators until done
	for {
		// exit if all test vals exhausted
		var done bool = true
		for _, gen := range colGens {
			if !gen.done() {
				done = false
				break
			}
		}
		if done {
			break
		}

		// insert data
		builder := newQueryBuilder(t.Namespace, t.Name)
		for _, gen := range colGens {
			builder.AddColumnData(gen.column.Name, gen.currentVal())
			gen.advance()
		}
		insertedVals, err := builder.RunQuery(conn)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if len(colGens) != len(insertedVals) {
			fmt.Printf("len mismatch inserted vals for table %s.%s\n", t.Namespace, t.Name)
			continue
		}

		// add inserted vals to record
		var crs []columnRecord
		for i := 0; i < len(colGens); i++ {
			cr := columnRecord{
				Name:   colGens[i].column.Name,
				Order:  colGens[i].column.Order,
				IsPKey: colGens[i].column.Constraint == query.CONSTRAINT_PRIMARY_KEY,
				Value:  insertedVals[i],
			}
			crs = append(crs, cr)
		}
		r.AddColumnRecords(t, crs)
	}
}
