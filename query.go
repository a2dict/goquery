package goquery

import (
	"math"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

func BuildPagedQuery(db *gorm.DB, entity interface{}) PagedQueryFunc {
	cols := parseDBCols(entity)
	colsFilter := ContainStringFilter(cols...)
	return func(qReq *QReq) (*PageWrap, error) {
		pw := &PageWrap{}
		out := reflect.New(reflect.SliceOf(reflect.TypeOf(entity))).Interface()
		preds := parseQueryMap(qReq.Q)

		limit := qReq.Size
		if limit <= 0 {
			limit = 20
		} else if limit > 200 {
			limit = 200
		}

		page := qReq.Page
		if page <= 0 {
			page = 1
		}

		var count int64
		q := db.Model(out)
		for _, p := range preds {
			if colsFilter(p.Col) {
				q = p.Apply(q)
			}
		}
		err := q.Count(&count).Error
		if err != nil {
			return nil, err
		}

		pages := int64(math.Ceil(float64(count) / float64(limit)))
		pw.Total = count
		pw.Size = limit
		pw.Page = qReq.Page
		pw.Pages = pages

		// sort
		for _, srt := range qReq.Sort {
			if col := strings.TrimLeft(srt, "+-"); colsFilter(col) {
				if strings.HasPrefix(srt, "-") {
					q = q.Order(col + " desc")
				} else {
					q = q.Order(col)
				}
			}
		}

		offset := (qReq.Page - 1) * limit
		q = q.Offset(offset).Limit(limit)

		err = q.Find(out).Error
		if err != nil {
			return nil, err
		}

		sliceOut := reflect.ValueOf(out).Elem().Interface()
		pw.Data = sliceOut
		return pw, nil
	}
}

func parseQueryMap(qm map[string]string) []*Predicate {
	var res []*Predicate
	for k, v := range qm {
		col, op := parseOp(k)
		cond := &Cond{Col: col, Op: op, Val: v}
		pred := cond.ToPredicate()
		if pred != nil {
			res = append(res, pred)
		}
	}
	return res
}

func parseDBCols(entity interface{}) []string {
	var res []string
	sc := gorm.Scope{Value: entity}
	ms := sc.GetModelStruct()
	for _, f := range ms.StructFields {
		if !f.IsIgnored {
			res = append(res, f.DBName)
		}
	}
	return res
}
