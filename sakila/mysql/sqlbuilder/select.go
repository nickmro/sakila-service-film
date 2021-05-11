package sqlbuilder

import (
	"fmt"
	"strings"
)

// SelectStatement builds a select query statement.
type SelectStatement struct {
	auxiliaries []*Auxiliary
	columns     []string
	table       string
	joins       []string
	conditions  []*Condition
	groups      []string
	sorts       []string
	limit       int
	offset      int
	unions      []*SelectStatement
	args        []interface{}
}

// SelectBuilder builds a select statement string.
type SelectBuilder struct {
	strings.Builder
	IndentSpaces int
}

// Select adds columns to the query.
func Select(columns ...string) *SelectStatement {
	return &SelectStatement{
		columns: columns,
	}
}

// With adds an auxiliary statement to the query.
func (s *SelectStatement) With(name string, as *SelectStatement) *SelectStatement {
	s.auxiliaries = append(s.auxiliaries, &Auxiliary{
		Name:   name,
		Select: as,
	})

	return s
}

// Columns adds columns to the query.
func (s *SelectStatement) Columns(columns ...string) *SelectStatement {
	s.columns = append(s.columns, columns...)
	return s
}

// From adds a table to the query.
func (s *SelectStatement) From(table string) *SelectStatement {
	s.table = table
	return s
}

// Where adds conditions to the query.
func (s *SelectStatement) Where(condition string, args ...interface{}) *SelectStatement {
	s.conditions = append(s.conditions, &Condition{Query: condition, Args: args})
	return s
}

// InnerJoin adds an inner join to the query.
func (s *SelectStatement) InnerJoin(join string) *SelectStatement {
	s.joins = append(s.joins, fmt.Sprintf("INNER JOIN %s", join))
	return s
}

// LeftJoin adds a left join to the query.
func (s *SelectStatement) LeftJoin(join string) *SelectStatement {
	s.joins = append(s.joins, fmt.Sprintf("LEFT JOIN %s", join))
	return s
}

// OrderBy adds sorts to the query.
func (s *SelectStatement) OrderBy(sorts ...string) *SelectStatement {
	s.sorts = append(s.sorts, sorts...)
	return s
}

// GroupBy adds groups to the query.
func (s *SelectStatement) GroupBy(groups ...string) *SelectStatement {
	s.groups = append(s.groups, groups...)
	return s
}

// Limit adds a limit to the query.
func (s *SelectStatement) Limit(limit int) *SelectStatement {
	s.limit = limit
	return s
}

// Offset adds an offset to the query.
func (s *SelectStatement) Offset(offset int) *SelectStatement {
	s.offset = offset
	return s
}

// UnionAll adds a union to the query.
func (s *SelectStatement) UnionAll(b *SelectStatement) *SelectStatement {
	s.unions = append(s.unions, b)
	return s
}

// Build builds the query and returns the args.
func (s *SelectStatement) Build() (string, []interface{}) {
	return s.String(), s.args
}

// String returns the select statement as a string.
func (s *SelectStatement) String() string {
	b := &SelectBuilder{}
	return b.Print(s)
}

// Pretty returns the select statement as a pretty printed string.
func (s *SelectStatement) Pretty() string {
	b := &SelectBuilder{}
	return b.PrettyPrint(s)
}

// Indent adds an indent to the statement print.
func (b *SelectBuilder) Indent(spaces int) *SelectBuilder {
	b.IndentSpaces = spaces
	return b
}

// Print prints a SQL statement.
func (b *SelectBuilder) Print(s *SelectStatement) string { //nolint
	if len(s.auxiliaries) > 0 {
		b.WriteString("WITH ")

		for i := range s.auxiliaries {
			if i > 0 {
				b.WriteString(", ")
			} else {
				b.WriteString(s.auxiliaries[i].String())
			}
		}

		b.WriteString(" ")
	}

	b.WriteString("SELECT ")

	for i := range s.columns {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(s.columns[i])
	}

	b.WriteString(" FROM ")
	b.WriteString(s.table)

	if len(s.joins) > 0 {
		b.WriteString(" ")

		for i := range s.joins {
			if i > 0 {
				b.WriteString(" ")
			}

			b.WriteString(s.joins[i])
		}
	}

	if len(s.conditions) > 0 {
		b.WriteString(" WHERE ")

		for i := range s.conditions {
			if i > 0 {
				b.WriteString(" AND ")
			}

			s.args = append(s.args, s.conditions[i].Args...)

			b.WriteString(s.conditions[i].String())
		}
	}

	if len(s.groups) > 0 {
		b.WriteString(" GROUP BY ")

		for i := range s.groups {
			if i > 0 {
				b.WriteString(", ")
			}

			b.WriteString(s.groups[i])
		}
	}

	if len(s.sorts) > 0 {
		b.WriteString(" ORDER BY ")

		for i := range s.sorts {
			if i > 0 {
				b.WriteString(", ")
			}

			b.WriteString(s.sorts[i])
		}
	}

	if s.limit > 0 {
		b.WriteString(" LIMIT ")
		b.WriteString(fmt.Sprintf("%d", s.limit))
	}

	if s.offset > 0 {
		b.WriteString(" OFFSET ")
		b.WriteString(fmt.Sprintf("%d", s.offset))
	}

	if len(s.unions) > 0 {
		for i := range s.unions {
			b.WriteString(" UNION ALL ")
			b.WriteString(s.unions[i].String())
		}
	}

	return b.String()
}

// PrettyPrint pretty prints a SQL statement.
func (b *SelectBuilder) PrettyPrint(s *SelectStatement) string { //nolint
	if len(s.auxiliaries) > 0 {
		b.WriteWithIndent("WITH ")

		withs := make([]string, len(s.auxiliaries))
		for i := range withs {
			withs[i] = s.auxiliaries[i].Pretty()
		}

		b.WriteString(strings.Join(withs, ",\n"+indent(b.IndentSpaces)))
		b.WriteString("\n")
	}

	b.WriteWithIndent("SELECT ")
	b.WriteString(strings.Join(s.columns, ",\n"+indent(b.IndentSpaces+7)))
	b.WriteString("\n")
	b.WriteWithIndent("FROM ")
	b.WriteString(s.table)

	if len(s.joins) > 0 {
		b.WriteString("\n")
		b.WriteWithIndent(strings.Join(s.joins, "\n"+indent(b.IndentSpaces)))
	}

	if len(s.conditions) > 0 {
		b.WriteString("\n")
		b.WriteWithIndent("WHERE ")

		for i := range s.conditions {
			if i > 0 {
				b.WriteString("\n" + indent(b.IndentSpaces+2) + "AND ")
			}

			b.WriteString(s.conditions[i].String())
		}
	}

	if len(s.groups) > 0 {
		b.WriteString("\n")
		b.WriteWithIndent("GROUP BY ")
		b.WriteString(strings.Join(s.groups, ",\n"+indent(b.IndentSpaces+9)))
	}

	if len(s.sorts) > 0 {
		b.WriteString("\n")
		b.WriteWithIndent("ORDER BY ")
		b.WriteString(strings.Join(s.sorts, ",\n"+indent(b.IndentSpaces+9)))
	}

	if s.limit > 0 {
		b.WriteString("\n")
		b.WriteWithIndent("LIMIT ")
		b.WriteString(fmt.Sprintf("%d", s.limit))
	}

	if s.limit > 0 {
		b.WriteString("\n")
		b.WriteWithIndent("OFFSET ")
		b.WriteString(fmt.Sprintf("%d", s.offset))
	}

	if len(s.unions) > 0 {
		for i := range s.unions {
			b.WriteString("\n")
			b.WriteWithIndent("UNION ALL")
			b.WriteString("\n")
			b.PrettyPrint(s.unions[i])
		}
	}

	return b.String()
}

// WriteWithIndent writes an indented string to the string builder.
func (b *SelectBuilder) WriteWithIndent(s string) {
	if b.IndentSpaces > 0 {
		b.WriteString(indent(b.IndentSpaces))
	}

	b.WriteString(s)
}

func indent(n int) string {
	return strings.Repeat(" ", n)
}
