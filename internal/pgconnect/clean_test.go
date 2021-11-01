package pgconnect_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/justasable/pgconnect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupDBWithScript(t *testing.T) {
	// connect to database
	conn, err := pgconnect.Connect()
	require.NoError(t, err)
	defer conn.Close(context.Background())

	// create table with data
	now := int(time.Now().Unix())
	testTable := "test" + strconv.Itoa(now)
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE public.%s (id integer)`, testTable))
	require.NoError(t, err)

	// run setup db with script
	conn.Close(context.Background())
	err = pgconnect.SetupDBWithScript("test_setup.sql")
	require.NoError(t, err)

	// reconnect to database
	conn, err = pgconnect.Connect()
	require.NoError(t, err)
	defer conn.Close(context.Background())

	// check if data has been wiped
	var exists bool
	err = conn.QueryRow(context.Background(), `
	SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = $1
	)`, testTable).Scan(&exists)
	assert.NoError(t, err)
	assert.Falsef(t, exists, "test table %s was not removed", testTable)

	// check setup script
	setupTable := "hello"
	err = conn.QueryRow(context.Background(), `
	SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE  table_schema = 'public'
		AND    table_name   = $1
	)`, setupTable).Scan(&exists)
	assert.NoError(t, err)
	assert.Truef(t, exists, "table %s from setup script not found", setupTable)
}
