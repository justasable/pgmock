package pgmock_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/justasable/pgmock/internal/pgmock"
	"github.com/stretchr/testify/assert"
)

const TEST_WIPE_TABLE = "testwipetable"

func TestDeleteData(t *testing.T) {
	// connect to database
	conn := MustConnect(t)
	defer conn.Close(context.Background())

	// check that the test table we are going to use does not exist
	var exists bool
	err := conn.QueryRow(context.Background(), `
	SELECT EXISTS (
		SELECT 1
		FROM pg_tables
		WHERE (schemaname='public' AND tablename=$1)
		AND (schemaname='test' AND tablename=$1)
	)`, TEST_WIPE_TABLE).Scan(&exists)
	assert.NoError(t, err)
	assert.Falsef(t, exists, "(public|test).%s should not exist, unable to run tests", TEST_WIPE_TABLE)

	// insert some data
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE public.%s (id integer)`, TEST_WIPE_TABLE))
	assert.NoError(t, err)
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE test.%s (id integer)`, TEST_WIPE_TABLE))
	assert.NoError(t, err)
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`INSERT INTO public.%s VALUES (1)`, TEST_WIPE_TABLE))
	assert.NoError(t, err)
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`INSERT INTO test.%s VALUES (1)`, TEST_WIPE_TABLE))
	assert.NoError(t, err)

	// run wipe data function
	err = pgmock.DeleteData(conn)
	assert.NoError(t, err)

	// check if data has been wiped
	var count int
	err = conn.QueryRow(
		context.Background(),
		fmt.Sprintf(`SELECT COUNT(*) FROM public.%s`, TEST_WIPE_TABLE),
	).Scan(&count)
	assert.NoError(t, err)
	assert.Zero(t, count)
	err = conn.QueryRow(
		context.Background(),
		fmt.Sprintf(`SELECT COUNT(*) FROM test.%s`, TEST_WIPE_TABLE),
	).Scan(&count)
	assert.NoError(t, err)
	assert.Zero(t, count)

	// remove tables
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`DROP TABLE public.%s`, TEST_WIPE_TABLE))
	assert.NoError(t, err)
	conn.Exec(context.Background(), fmt.Sprintf(`DROP TABLE test.%s`, TEST_WIPE_TABLE))
	assert.NoError(t, err)
}
