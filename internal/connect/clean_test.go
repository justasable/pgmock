package connect_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/justasable/pgmock/internal/connect"
	"github.com/stretchr/testify/assert"
)

func TestDropAndRecreateDB(t *testing.T) {
	// connect to database
	conn, err := connect.Connect()
	assert.NoError(t, err)
	defer conn.Close(context.Background())

	// create table with data
	now := int(time.Now().Unix())
	testTable := "test" + strconv.Itoa(now)
	_, err = conn.Exec(context.Background(), fmt.Sprintf(`CREATE TABLE public.%s (id integer)`, testTable))
	assert.NoError(t, err)

	// run wipe data function
	conn.Close(context.Background())
	err = connect.DropAndRecreateDB()
	assert.NoError(t, err)

	// reconnect to database
	conn, err = connect.Connect()
	assert.NoError(t, err)
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
	assert.False(t, exists)
}

func TestRunScript(t *testing.T) {
	err := connect.RunScript("test_script.sql")
	assert.NoError(t, err)

	err = connect.RunScript("non_existent_script.sql")
	assert.Error(t, err)
}
