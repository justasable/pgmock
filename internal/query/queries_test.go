package query_test

import (
	"context"
	"testing"

	"github.com/justasable/pgmock/internal/pgconnect"
	"github.com/justasable/pgmock/internal/query"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/suite"
)

type QuerySuite struct {
	suite.Suite
	conn *pgx.Conn
}

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(QuerySuite))
}

func (s *QuerySuite) SetupTest() {
	config, err := pgx.ParseConfig("")
	s.NoError(err)
	err = pgconnect.SetupDBWithScript(config, "test_setup.sql")
	s.NoError(err)

	conn, err := pgx.Connect(context.Background(), "")
	s.NoError(err)

	s.conn = conn
}

func (s *QuerySuite) TearDownTest() {
	s.conn.Close(context.Background())
}

func (s *QuerySuite) TestTables() {
	// fetch tables
	tt, err := query.Tables(s.conn)
	s.Require().NoError(err)

	// expected tables
	expected := []query.Table{
		{Namespace: "public", Name: "types", Columns: []query.Column{
			{Order: 1, Name: "type_int", DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 2, Name: "type_bool", DataType: pgtype.BoolOID, TypeMod: -1},
			{Order: 3, Name: "type_numeric", DataType: pgtype.NumericOID, TypeMod: -1},
			{Order: 4, Name: "type_numericp", DataType: pgtype.NumericOID, TypeMod: 655366},
			{Order: 5, Name: "type_text", DataType: pgtype.TextOID, TypeMod: -1},
			{Order: 6, Name: "type_timestamptz", DataType: pgtype.TimestamptzOID, TypeMod: -1},
			{Order: 7, Name: "type_date", DataType: pgtype.DateOID, TypeMod: -1},
			{Order: 8, Name: "type_byte", DataType: pgtype.ByteaOID, TypeMod: -1},
			{Order: 9, Name: "type_uuid", DataType: pgtype.UUIDOID, TypeMod: -1},
		}},
		{Namespace: "public", Name: "identity", Columns: []query.Column{
			{Order: 1, Name: "identity_always", IsNotNull: true, Identity: query.IDENTITY_ALWAYS, Constraint: query.CONSTRAINT_PRIMARY_KEY, ConKeys: []int{1}, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 2, Name: "identity_default", IsNotNull: true, Identity: query.IDENTITY_DEFAULT, DataType: pgtype.Int4OID, TypeMod: -1},
		}},
		{Namespace: "public", Name: "constraints", Columns: []query.Column{
			{Order: 1, Name: "con_pk_one", IsNotNull: true, Constraint: query.CONSTRAINT_PRIMARY_KEY, ConKeys: []int{1, 2}, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 2, Name: "con_pk_two", IsNotNull: true, Constraint: query.CONSTRAINT_PRIMARY_KEY, ConKeys: []int{1, 2}, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 3, Name: "con_null", DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 4, Name: "con_not_null", IsNotNull: true, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 5, Name: "con_default", HasDefault: true, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 6, Name: "con_no_default", DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 7, Name: "con_generated", HasDefault: true, Generated: query.GENERATED_STORED, DataType: pgtype.Int4OID, TypeMod: -1},
		}},
		{Namespace: "test", Name: "references", Columns: []query.Column{
			{Order: 1, Name: "fk_single", Constraint: query.CONSTRAINT_FOREIGN_KEY, ConKeys: []int{1}, FKColumns: []int{1}, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 2, Name: "fk_multiple_one", Constraint: query.CONSTRAINT_FOREIGN_KEY, ConKeys: []int{2, 3}, FKColumns: []int{1, 2}, DataType: pgtype.Int4OID, TypeMod: -1},
			{Order: 3, Name: "fk_multiple_two", Constraint: query.CONSTRAINT_FOREIGN_KEY, ConKeys: []int{2, 3}, FKColumns: []int{1, 2}, DataType: pgtype.Int4OID, TypeMod: -1},
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
		s.Require().FailNowf("tables not found", "table: %s.%s", expected[i].Namespace, expected[i].Name)
	}

	// fill in foreign key oids
	expected[3].Columns[0].FKTableID = expected[1].ID
	expected[3].Columns[1].FKTableID = expected[2].ID
	expected[3].Columns[2].FKTableID = expected[2].ID

	// compare expected with got
	for _, e := range expected {
		for _, aTable := range tt {
			if e.ID == aTable.ID {
				s.Equalf(e, aTable, "mismatch in table: %s.%s", e.Namespace, e.Name)
			}
		}
	}
}
