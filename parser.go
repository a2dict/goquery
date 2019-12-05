package goquery

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/jinzhu/gorm"
)

// verb
const (
	// GT ...
	GT = "gt"
	// GE ...
	GE = "ge"
	// LT ...
	LT = "lt"
	// LE ...
	LE = "le"
	// IN ...
	IN = "in"
	// NN ...
	NN = "not_in"
	// LIKE ...
	LIKE = "like"
	// ILIKE ...
	ILIKE = "ilike"
	// EQ ...
	EQ = "eq"
)

// Apply ...
func (p *Predicate) Apply(q *gorm.DB) *gorm.DB {
	return q.Where(p.Where, p.Val)
}

// Cond ...
type Cond struct {
	Op  string
	Col string
	Val string
}

// Predicate ...
type Predicate struct {
	Col   string
	Where string
	Val   interface{}
}

// ToPredicate ...
func (c *Cond) ToPredicate() *Predicate {
	switch c.Op {
	case GT:
		return &Predicate{Where: fmt.Sprintf("`%v` > ?", c.Col), Val: c.Val, Col: c.Col}
	case GE:
		return &Predicate{Where: fmt.Sprintf("`%v` >= ?", c.Col), Val: c.Val, Col: c.Col}
	case LT:
		return &Predicate{Where: fmt.Sprintf("`%v` < ?", c.Col), Val: c.Val, Col: c.Col}
	case LE:
		return &Predicate{Where: fmt.Sprintf("`%v` <= ?", c.Col), Val: c.Val, Col: c.Col}
	case IN:
		var v []interface{}
		json.Unmarshal([]byte(c.Val), &v)
		return &Predicate{Where: fmt.Sprintf("`%v` in (?)", c.Col), Val: v, Col: c.Col}
	case NN:
		var v []interface{}
		json.Unmarshal([]byte(c.Val), &v)
		return &Predicate{Where: fmt.Sprintf("`%v` not in (?)", c.Col), Val: v, Col: c.Col}
	case LIKE:
		return &Predicate{Where: fmt.Sprintf("`%v` like ?", c.Col), Val: c.Val, Col: c.Col}
	case ILIKE:
		return &Predicate{Where: fmt.Sprintf("lower(`%v`) like lower(?)", c.Col), Val: c.Val, Col: c.Col}
	case EQ:
		return &Predicate{Where: fmt.Sprintf("`%v` = ?", c.Col), Val: c.Val, Col: c.Col}
	default:
		return nil
	}
}

func parseOp(s string) (col, op string) {
	exp := regexp.MustCompile("^(\\w+)(::(\\w+))?$")
	sm := exp.FindStringSubmatch(s)
	col = sm[1]
	op = EQ
	if len(sm) >= 4 && sm[3] != "" {
		op = sm[3]
	}
	return
}
