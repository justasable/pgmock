package pgmock

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// fetchAllTables retrieves tables from all available non-system namespaces
// filters with following rules:
// 	- table only (i.e. no views)
// 	- excludes any other session's temporary tables
// 	- excludes default system namespaces
// ordered by table name in alphabetical ascending
func FetchAllTables(conn *pgx.Conn) ([]Table, error) {
	// query tables
	rows, err := conn.Query(
		context.Background(),
		`SELECT pgc.oid, nsp.nspname, pgc.relname
		FROM pg_class pgc
			JOIN pg_namespace nsp ON pgc.relnamespace = nsp.oid
		WHERE 
			pgc.relkind = ANY (ARRAY['r'::"char", 'p'::"char"])
			AND NOT pg_is_other_temp_schema(nsp.oid)
			AND NOT nsp.nspname = ANY(ARRAY['pg_catalog', 'pg_toast', 'information_schema'])
		ORDER BY pgc.relname ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// scan tables results
	var tt []Table
	for rows.Next() {
		var t Table
		if err = rows.Scan(&t.ID, &t.Namespace, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan results: %w", err)
		}

		tt = append(tt, t)
	}

	// fetch columns for tables
	var tables []Table
	for _, aTable := range tt {
		t := aTable
		cols, err := fetchColumns(conn, aTable.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch columns: %w", err)
		}

		t.Columns = cols
		tables = append(tables, t)
	}

	return tables, nil
}

// fetchColumns grabs all 'relevant' columns of a specified table
// relevant is defined as non-generated, non-always identity columns
func fetchColumns(conn *pgx.Conn, tableID pgtype.OID) ([]Column, error) {
	// query table columns
	rows, err := conn.Query(
		context.Background(),
		`SELECT
			att.attnum, att.attname, NOT att.attnotnull, att.atthasdef, att.attidentity='d',
			COALESCE(con.contype='p', FALSE), COALESCE(con.confrelid, 0), con.confkey
		FROM pg_attribute att
		LEFT OUTER JOIN pg_constraint con
			ON att.attrelid = con.conrelid AND att.attnum = ANY(con.conkey)
		WHERE att.attrelid=$1
			AND att.attnum > 0
			AND NOT att.attisdropped
			AND att.attidentity!='a'
			AND att.attgenerated=''
		ORDER BY att.attnum ASC
	`, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// scan column results
	var cols []Column
	for rows.Next() {
		var c Column
		if err = rows.Scan(
			&c.Order, &c.Name, &c.IsNullable, &c.HasDefaultValue,
			&c.IsDefaultIdentity, &c.IsPrimaryKey, &c.FKTableID,
			&c.FKColumns,
		); err != nil {
			return nil, err
		}

		cols = append(cols, c)
	}

	return cols, nil
}
