package pgmock_test

import (
	"fmt"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/pgmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllTables(t *testing.T) {
	// fetch tables
	conn := MustConnect(t)
	tt, err := pgmock.FetchAllTables(conn)
	assert.NoError(t, err)

	// check number of tables returned
	assert.Len(t, tt, 3)

	// grab dynamically generated table ids
	var tablePublicAccount pgtype.OID
	var tablePublicPet pgtype.OID
	var tableTestAccount pgtype.OID
	for _, aTable := range tt {
		if aTable.Namespace == "public" && aTable.Name == "account" {
			tablePublicAccount = aTable.ID
		} else if aTable.Namespace == "public" && aTable.Name == "pet" {
			tablePublicPet = aTable.ID
		} else if aTable.Namespace == "test" && aTable.Name == "account" {
			tableTestAccount = aTable.ID
		}
	}

	// compare fetched tables
	expected := []pgmock.Table{
		{ID: tablePublicAccount, Namespace: "public", Name: "account", Columns: []pgmock.Column{
			{Order: 2, Name: "is_not_nullable"},
			{Order: 3, Name: "is_nullable", IsNullable: true},
		}},
		{ID: tablePublicPet, Namespace: "public", Name: "pet", Columns: []pgmock.Column{
			{Order: 1, Name: "fname", IsPrimaryKey: true},
			{Order: 2, Name: "lname", IsPrimaryKey: true},
		}},
		{ID: tableTestAccount, Namespace: "test", Name: "account", Columns: []pgmock.Column{
			{Order: 1, Name: "id_default", IsDefaultIdentity: true, IsPrimaryKey: true},
			{Order: 2, Name: "text_default", IsNullable: true, HasDefaultValue: true},
			{Order: 3, Name: "fk_single", IsNullable: true, FKTableID: tablePublicAccount, FKColumns: []int{1}},
			{Order: 4, Name: "fk_multiple_fname", IsNullable: true, FKTableID: tablePublicPet, FKColumns: []int{1, 2}},
			{Order: 5, Name: "fk_multiple_lname", IsNullable: true, FKTableID: tablePublicPet, FKColumns: []int{1, 2}},
		}},
	}
	for _, e := range expected {
		assert.Containsf(t, tt, e, fmt.Sprintf("Testing table: %s.%s", e.Namespace, e.Name))
	}
}
