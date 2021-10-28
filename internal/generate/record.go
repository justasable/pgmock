package generate

import (
	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/query"
)

type recordSet map[pgtype.OID]tableRecord

func (r recordSet) HasData(tableID pgtype.OID) bool {
	_, ok := r[tableID]
	return ok
}

func (r recordSet) AddColumnRecords(t query.Table, records []columnRecord) {
	// create table record entry if needed
	if _, ok := r[t.ID]; !ok {
		r[t.ID] = tableRecord{
			tableID: t.ID,
			schema:  t.Namespace,
			name:    t.Name,
		}
	}

	// add to table records
	tr := r[t.ID]
	tr.rowRecords = append(tr.rowRecords, rowRecord{columnRecords: records})
	r[t.ID] = tr
}

func (r recordSet) PKeyForTable(tableID pgtype.OID) []columnRecord {
	tr, ok := r[tableID]
	if !ok || len(tr.rowRecords) == 0 {
		return nil
	}

	var pkeys []columnRecord
	for _, cr := range tr.rowRecords[0].columnRecords {
		if cr.IsPKey {
			pkeys = append(pkeys, cr)
		}
	}

	return pkeys
}

func (r recordSet) TotalRecords() int {
	var count int
	for _, tr := range r {
		count += len(tr.rowRecords)
	}
	return count
}

type tableRecord struct {
	tableID    pgtype.OID
	schema     string
	name       string
	rowRecords []rowRecord
}

type rowRecord struct {
	columnRecords []columnRecord
}

type columnRecord struct {
	Name   string
	Order  int
	IsPKey bool
	Value  interface{}
}
