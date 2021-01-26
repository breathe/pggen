package pginfer

import (
	"github.com/jschaf/pggen/internal/ast"
	"github.com/jschaf/pggen/internal/pg"
	"strings"
	"unicode"
)

// isColNullable tries to prove the column is not nullable. Strive for
// correctness here: it's better to assume a column is nullable when we can't
// know for sure.
func isColNullable(query *ast.SourceQuery, plan Plan, out string, column pg.Column) bool {
	switch {
	case len(out) == 0:
		// No output? Not sure what this means but do the check here so we don't
		// have to do it in each case below.
		return false
	case strings.HasPrefix(out, "'"):
		return false // literal string can't be null
	case unicode.IsDigit(rune(out[0])):
		return false // literal number can't be null
	default:
		// try below
	}

	if plan.Type == PlanModifyTable && plan.Relation == column.TableName && !column.Null {
		// A returning clause in an insert, update, or delete statement. The column
		// must come from the underlying table and must have a not null constraint.
		return false
	}
	return true // we can't figure it out; assume nullable
}
