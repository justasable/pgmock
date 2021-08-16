package pgmock

import (
	"context"
	"fmt"

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

	// scan results and fetch table columns
	var tables []Table
	for rows.Next() {
		var t Table
		if err = rows.Scan(&t.ID, &t.Namespace, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan results: %w", err)
		}

		tables = append(tables, t)
	}

	return tables, nil
}
