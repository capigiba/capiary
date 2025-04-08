package query

import (
	"fmt"
	"strings"
)

// BuildPostgresSelectQuery builds a SELECT query and its arguments for PostgreSQL.
func BuildPostgresSelectQuery(tableName string, opts QueryOptions) (string, []interface{}) {
	var (
		args         []interface{}
		whereClauses []string
		selectFields string
		orderClause  string
	)

	// If no specific fields, select "*"
	if len(opts.Fields) == 0 {
		selectFields = "*"
	} else {
		selectFields = strings.Join(opts.Fields, ", ")
	}

	// Build WHERE clauses
	argIndex := 1
	for _, fil := range opts.Filters {
		var op string
		switch fil.Operator {
		case OpEqual:
			op = "="
		case OpNotEqual:
			op = "!="
		case OpGreaterThan:
			op = ">"
		case OpLessThan:
			op = "<"
		case OpGTE:
			op = ">="
		case OpLTE:
			op = "<="
		default:
			// Fallback or handle custom operators here
			op = "="
		}

		whereClauses = append(whereClauses, fmt.Sprintf("%s %s $%d", fil.Field, op, argIndex))
		args = append(args, fil.Value)
		argIndex++
	}

	// Build ORDER BY
	if len(opts.Sorts) > 0 {
		var sortExprs []string
		for _, s := range opts.Sorts {
			dir := "ASC"
			if s.Desc {
				dir = "DESC"
			}
			sortExprs = append(sortExprs, fmt.Sprintf("%s %s", s.Field, dir))
		}
		orderClause = "ORDER BY " + strings.Join(sortExprs, ", ")
	}

	// Final SQL
	query := fmt.Sprintf("SELECT %s FROM %s", selectFields, tableName)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	if orderClause != "" {
		query += " " + orderClause
	}

	return query, args
}
