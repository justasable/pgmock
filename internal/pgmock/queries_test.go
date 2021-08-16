package pgmock_test

import (
	"testing"

	"github.com/justasable/pgmock/internal/pgmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllTables(t *testing.T) {
	// fetch tables
	conn := MustConnect(t)
	tt, err := pgmock.FetchAllTables(conn)
	assert.NoError(t, err)

	// check number of tables returned
	assert.Len(t, tt, 2)

	// set table oid to 0 as the actual oid is unknown
	got := tt[:0]
	for _, t := range tt {
		t.ID = 0
		got = append(got, t)
	}

	// check against expected test tables
	expected := []pgmock.Table{
		{Namespace: "public", Name: "account"},
		{Namespace: "test", Name: "account"},
	}

	assert.ElementsMatch(t, got, expected)
}
