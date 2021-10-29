// package pgtestvals generates simple postgresql data type test values and unique values
package pgtestvals

import "github.com/jackc/pgtype"

func TestVals(dataType pgtype.OID, typeMod int) []interface{} {
	switch dataType {
	case pgtype.Int4OID:
		return integerTestVals
	case pgtype.BoolOID:
		return boolTestVals
	case pgtype.NumericOID:
		return numericTestVals(typeMod)
	case pgtype.TextOID:
		return textTestVals
	case pgtype.TimestamptzOID:
		return timestatmptzTestVals
	case pgtype.DateOID:
		return dateTestVals
	case pgtype.ByteaOID:
		return byteTestVals
	case pgtype.UUIDOID:
		return uuidTestVals
	}

	return nil
}

func UniqueFn(dataType pgtype.OID, typeMod int) func(int) interface{} {
	switch dataType {
	case pgtype.Int4OID:
		return integerUniqueFn
	case pgtype.BoolOID:
		return boolUniqueFn
	case pgtype.NumericOID:
		return numericUniqueFn(typeMod)
	case pgtype.TextOID:
		return textUniqueFn
	case pgtype.TimestamptzOID:
		return timestamptzUniqueFn
	case pgtype.DateOID:
		return dateUniqueFn
	case pgtype.ByteaOID:
		return byteUniqueFn
	case pgtype.UUIDOID:
		return uuidUniqueFn
	}

	return nil
}
