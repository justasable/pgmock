// package generate provides fake values of specified types
package generate

// // GenerateData generates mock data for all given tables in the db
// func GenerateData(conn *pgx.Conn) error {
// 	tt, err := query.Tables(conn)
// 	if err != nil {
// 		return err
// 	}

// 	// first pass: generate rows for column types
// 	for _, t := range tt {
// 		err := generateForTableValues(conn, &t)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func generateForTableValues(conn *pgx.Conn, table *query.Table) error {
// 	// create task for each column
// 	var tt []*colTask
// 	for i := 0; i < len(table.Columns); i++ {
// 		t, err := newColTask(table.Columns[i])
// 		if err != nil {
// 			return err
// 		} else if t == nil {
// 			continue
// 		}
// 		tt = append(tt, t)
// 	}

// 	// generate data row until all tasks finished
// 	var errs []error
// 	for {
// 		// check if all tasks finished
// 		var finished bool = true
// 		for _, t := range tt {
// 			if !t.Done() {
// 				finished = false
// 			}
// 		}
// 		if finished {
// 			break
// 		}

// 		// gather values for row
// 		var colNames []string
// 		var values []interface{}
// 		var placeholders []string
// 		var defaultValCount int
// 		for i, t := range tt {
// 			colNames = append(colNames, t.column.Name)
// 			if _, ok := t.CurrentVal().(defaultType); ok {
// 				placeholders = append(placeholders, "DEFAULT")
// 				defaultValCount++
// 			} else {
// 				placeholders = append(placeholders, fmt.Sprintf("$%d", i+1-defaultValCount))
// 				values = append(values, t.CurrentVal())
// 			}

// 			t.Advance()
// 		}

// 		// build db command
// 		cmd := fmt.Sprintf(
// 			"INSERT INTO %s.%s (%s) VALUES (%s)",
// 			table.Namespace, table.Name,
// 			strings.Join(colNames, ","),
// 			strings.Join(placeholders, ","),
// 		)

// 		// insert row
// 		_, err := conn.Exec(context.Background(), cmd, values...)
// 		if err != nil {
// 			// we don't return on this error as some unforseen constraints may prevent
// 			// row insert but we can still continue with the next row
// 			errs = append(errs, err)
// 		}
// 	}

// 	if len(errs) > 0 {
// 		var errStr string
// 		for _, err := range errs {
// 			errStr += fmt.Sprintf("%s\n", err)
// 		}
// 		return errors.New(errStr)
// 	}

// 	return nil
// }
