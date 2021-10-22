package generate

import (
	"errors"

	"github.com/jackc/pgtype"
)

var ErrUnsupportedType = errors.New("unsupported data type")

type Generator struct {
	defaultVals []interface{}
	uniqueGen   func(n int) interface{}
}

func NewGenerator(oid pgtype.OID) (Generator, error) {
	g := Generator{}
	switch oid {
	case pgtype.Int4OID:
		for _, val := range IntegerDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return IntegerUnique(n) }
	case pgtype.BoolOID:
		for _, val := range BooleanDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return BooleanUnique(n) }
	case pgtype.NumericOID:
		for _, val := range NumericDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return NumericUnique(n) }
	case pgtype.TextOID:
		for _, val := range TextDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return TextUnique(n) }
	case pgtype.TimestamptzOID:
		for _, val := range TimestampTZDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return TimestamptTZUnique(n) }
	case pgtype.DateOID:
		for _, val := range DateDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return DateUnique(n) }
	case pgtype.ByteaOID:
		for _, val := range ByteDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return ByteUnique(n) }
	case pgtype.UUIDOID:
		for _, val := range UUIDDefaults() {
			g.defaultVals = append(g.defaultVals, val)
		}
		g.uniqueGen = func(n int) interface{} { return UUIDUnique(n) }
	default:
		return g, ErrUnsupportedType
	}

	return g, nil
}

func (g Generator) ValueForRow(n int) interface{} {
	if n < len(g.defaultVals) {
		return g.defaultVals[n]
	}

	return g.uniqueGen(n - len(g.defaultVals))
}

func (g Generator) Done(n int) bool {
	return n >= len(g.defaultVals)
}
