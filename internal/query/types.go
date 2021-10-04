package query

import (
	"github.com/jackc/pgtype"
)

type Table struct {
	ID        pgtype.OID
	Namespace string
	Name      string
	Columns   []Column
}

type Column struct {
	Order             int
	Name              string
	IsNullable        bool
	HasDefaultValue   bool
	IsDefaultIdentity bool
	IsPrimaryKey      bool
	FKTableID         pgtype.OID
	FKColumns         []int
}
