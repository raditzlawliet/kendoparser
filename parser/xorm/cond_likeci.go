package kpxorm

import (
	"fmt"

	"xorm.io/builder"
)

// LikeCi defines likeci condition
type LikeCi [2]string

var _ builder.Cond = LikeCi{"", ""}

// WriteTo write SQL to Writer
func (likeci LikeCi) WriteTo(w builder.Writer) error {
	if _, err := fmt.Fprintf(w, "LOWER(%s) LIKE LOWER(?)", likeci[0]); err != nil {
		return err
	}
	// FIXME: if use other regular express, this will be failed. but for compatible, keep this
	if likeci[1][0] == '%' || likeci[1][len(likeci[1])-1] == '%' {
		w.Append(likeci[1])
	} else {
		w.Append("%" + likeci[1] + "%")
	}
	return nil
}

// And implements And with other conditions
func (likeci LikeCi) And(conds ...builder.Cond) builder.Cond {
	return builder.And(likeci, builder.And(conds...))
}

// Or implements Or with other conditions
func (likeci LikeCi) Or(conds ...builder.Cond) builder.Cond {
	return builder.Or(likeci, builder.Or(conds...))
}

// IsValid tests if this condition is valid
func (likeci LikeCi) IsValid() bool {
	return len(likeci[0]) > 0 && len(likeci[1]) > 0
}
