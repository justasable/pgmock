package query

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// tables retrieves tables from all available non-system namespaces
// filters with following rules:
// 	- table only (i.e. no views)
// 	- excludes any other session's temporary tables
// 	- excludes default system namespaces
func Tables(conn *pgx.Conn) ([]Table, error) {
	tt, err := tables(conn)
	if err != nil {
		return nil, err
	}

	var ret []Table
	for _, t := range tt {
		cc, err := columns(conn, t.ID)
		if err != nil {
			return nil, err
		}
		aTable := t
		aTable.Columns = cc
		ret = append(ret, aTable)
	}

	return ret, nil
}

func tables(conn *pgx.Conn) ([]Table, error) {
	// db query
	rows, err := conn.Query(
		context.Background(),
		`SELECT pgc.oid, nsp.nspname, pgc.relname
		FROM pg_class pgc
			JOIN pg_namespace nsp ON pgc.relnamespace = nsp.oid
		WHERE 
			pgc.relkind = ANY (ARRAY['r'::"char", 'p'::"char"])
			AND NOT pg_is_other_temp_schema(nsp.oid)
			AND NOT nsp.nspname = ANY(ARRAY['pg_catalog', 'pg_toast', 'information_schema'])
		ORDER BY nsp.nspname, pgc.relname`)
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

	return tt, nil
}

func columns(conn *pgx.Conn, tableID pgtype.OID) ([]Column, error) {
	// db query
	rows, err := conn.Query(
		context.Background(),
		`SELECT
			att.attnum, att.attname, att.attnotnull, att.atthasdef, att.attidentity,
			att.attgenerated, COALESCE(con.contype, ''), con.conkey,
			COALESCE(con.confrelid, 0), con.confkey, att.atttypid, att.atttypmod
		FROM pg_attribute att
		LEFT OUTER JOIN pg_constraint con
			ON att.attrelid = con.conrelid AND att.attnum = ANY(con.conkey)
		WHERE att.attrelid=$1
			AND att.attnum > 0
			AND NOT att.attisdropped
		ORDER BY att.attnum`,
		tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// scan column results
	var cols []Column
	for rows.Next() {
		var c Column
		if err = rows.Scan(
			&c.Order, &c.Name, &c.IsNotNull, &c.HasDefault, &c.Identity,
			&c.Generated, &c.Constraint, &c.ConKeys,
			&c.FKTableID, &c.FKColumns, &c.DataType, &c.TypeMod,
		); err != nil {
			return nil, err
		}

		cols = append(cols, c)
	}

	return cols, nil
}
