package kpxorm

import (
	"fmt"
	"sort"

	"xorm.io/builder"
)

// EqCi defines equal conditions
type EqCi map[string]interface{}

var _ builder.Cond = EqCi{}

// OpWriteTo writes conditions with special operator
func (eqci EqCi) OpWriteTo(op string, w builder.Writer) error {
	var i = 0
	for _, k := range eqci.sortedKeys() {
		v := eqci[k]
		switch v.(type) {
		case []int, []int64, []string, []int32, []int16, []int8, []uint, []uint64, []uint32, []uint16, []interface{}:
			if err := builder.In(k, v).WriteTo(w); err != nil {
				return err
			}
		case *builder.Builder:
			if _, err := fmt.Fprintf(w, "LOWER(%s)=LOWER(", k); err != nil {
				return err
			}

			if err := v.(*builder.Builder).WriteTo(w); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		case nil:
			if _, err := fmt.Fprintf(w, "%s=null", k); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "LOWER(%s)=LOWER(?)", k); err != nil {
				return err
			}
			w.Append(v)
		}
		if i != len(eqci)-1 {
			if _, err := fmt.Fprint(w, op); err != nil {
				return err
			}
		}
		i = i + 1
	}
	return nil
}

// WriteTo writes SQL to Writer
func (eqci EqCi) WriteTo(w builder.Writer) error {
	return eqci.OpWriteTo(" AND ", w)
}

// And implements And with other conditions
func (eqci EqCi) And(conds ...builder.Cond) builder.Cond {
	return builder.And(eqci, builder.And(conds...))
}

// Or implements Or with other conditions
func (eqci EqCi) Or(conds ...builder.Cond) builder.Cond {
	return builder.Or(eqci, builder.Or(conds...))
}

// IsValid tests if this EqCi is valid
func (eqci EqCi) IsValid() bool {
	return len(eqci) > 0
}

// sortedKeys returns all keys of this EqCi sorted with sort.Strings.
// It is used internally for consistent ordering when generating
// SQL, see https://gitea.com/xorm/builder/issues/10
func (eqci EqCi) sortedKeys() []string {
	keys := make([]string, 0, len(eqci))
	for key := range eqci {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
