package pgmock

import (
	"context"

	"github.com/jackc/pgx/v4"
)

// DeleteData removes data from all user tables leaving only schema
func DeleteData(conn *pgx.Conn) error {
	// remove data from all user defined tables
	conn.Exec(context.Background(), `
	DO $$
	BEGIN
		EXECUTE
		(SELECT 'TRUNCATE TABLE '
			|| string_agg(format('%I.%I', schemaname, tablename), ', ')
			|| ' CASCADE'
		FROM   pg_tables
		WHERE schemaname != 'pg_catalog'
			AND schemaname != 'information_schema'
			AND schemaname != 'pg_toast'
		);
	END$$;`)

	return nil
}
