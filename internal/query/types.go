package query

import "github.com/jackc/pgtype"

type Table struct {
	ID        pgtype.OID
	Namespace string
	Name      string
	Columns   []Column
}

type Column struct {
	Order      int
	Name       string
	IsNotNull  bool
	HasDefault bool
	Identity   Identity
	Generated  Generated
	Constraint Constraint
	FKTableID  pgtype.OID
	FKColumns  []int
	DataType   pgtype.OID
}

type Identity int8

const (
	IDENTITY_NONE    Identity = 0
	IDENTITY_ALWAYS  Identity = 97  // 'a'
	IDENTITY_DEFAULT Identity = 100 // 'd'
)

type Generated int8

const (
	GENERATED_NONE   Generated = 0
	GENERATED_STORED Generated = 115 // 's'
)

type Constraint int8

const (
	CONSTRAINT_NONE        Constraint = 0
	CONSTRAINT_PRIMARY_KEY Constraint = 112 // 'p'
	CONSTRAINT_FOREIGN_KEY Constraint = 102 // 'f'
	CONSTRAINT_CHECK       Constraint = 99  // 'c'
	CONSTRAINT_UNIQUE      Constraint = 117 // 'u'
	CONSTRAINT_TRIGGER     Constraint = 116 // 't'
	CONSTRAINT_EXCLUSION   Constraint = 120 // 'x'
)
