package query_test

import (
	"github.com/justasable/pgmock/internal/query"
)

func (s *DefaultSuite) TestTables() {
	// fetch tables
	tt, err := query.Tables(s.conn)
	s.Require().NoError(err)

	// expected tables
	expected := []query.Table{
		{Namespace: "public", Name: "types", Columns: []query.Column{
			{Order: 1, Name: "id", IsNotNull: true, Identity: query.IDENTITY_ALWAYS, Constraint: query.CONSTRAINT_PRIMARY_KEY},
			{Order: 2, Name: "id_default", IsNotNull: true, Identity: query.IDENTITY_DEFAULT},
		}},
		{Namespace: "public", Name: "constraints", Columns: []query.Column{
			{Order: 1, Name: "con_pk_one", IsNotNull: true, Constraint: query.CONSTRAINT_PRIMARY_KEY},
			{Order: 2, Name: "con_pk_two", IsNotNull: true, Constraint: query.CONSTRAINT_PRIMARY_KEY},
			{Order: 3, Name: "con_null"},
			{Order: 4, Name: "con_not_null", IsNotNull: true},
			{Order: 5, Name: "con_default", HasDefault: true},
			{Order: 6, Name: "con_no_default"},
			{Order: 7, Name: "con_generated", HasDefault: true, Generated: query.GENERATED_STORED},
		}},
		{Namespace: "test", Name: "references", Columns: []query.Column{
			{Order: 1, Name: "fk_single", Constraint: query.CONSTRAINT_FOREIGN_KEY, FKColumns: []int{1}},
			{Order: 2, Name: "fk_multiple_one", Constraint: query.CONSTRAINT_FOREIGN_KEY, FKColumns: []int{1, 2}},
			{Order: 3, Name: "fk_multiple_two", Constraint: query.CONSTRAINT_FOREIGN_KEY, FKColumns: []int{1, 2}},
		}},
	}

	// check number of tables returned
	s.Require().Len(tt, len(expected))

	// fill in table oids
out:
	for i := 0; i < len(expected); i++ {
		for _, aTable := range tt {
			if expected[i].Namespace == aTable.Namespace && expected[i].Name == aTable.Name {
				expected[i].ID = aTable.ID
				continue out
			}
		}
		s.Failf("table not found: %s.%s", expected[i].Namespace, expected[i].Name)
	}

	// fill in foreign key oids
	expected[2].Columns[0].FKTableID = expected[0].ID
	expected[2].Columns[1].FKTableID = expected[1].ID
	expected[2].Columns[2].FKTableID = expected[1].ID

	// compare expected with got
	for _, e := range expected {
		for _, aTable := range tt {
			if e.ID == aTable.ID {
				s.Equalf(e, aTable, "mismatch in table: %s.%s", e.Namespace, e.Name)
			}
		}
	}
}
