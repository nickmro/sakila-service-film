package sqlbuilder

import "strings"

// Auxiliary is a SQL auxiliary statement.
type Auxiliary struct {
	Name   string
	Select *SelectStatement
}

const (
	auxiliaryIndentSpaces = 4
)

// String returns the auxiliary statement as a string.
func (a *Auxiliary) String() string {
	b := strings.Builder{}
	b.WriteString(a.Name)
	b.WriteString(" AS (")
	b.WriteString(a.Select.String())
	b.WriteString(")")

	return b.String()
}

// Pretty returns a pretty printed auxiliary statement string.
func (a *Auxiliary) Pretty() string {
	b := strings.Builder{}

	b.WriteString(a.Name)
	b.WriteString(" AS (\n")

	s := &SelectBuilder{IndentSpaces: auxiliaryIndentSpaces}
	b.WriteString(s.PrettyPrint(a.Select))

	b.WriteString("\n)")

	return b.String()
}
