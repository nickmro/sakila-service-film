package sqlbuilder

import "fmt"

// Condition represents a where clause.
type Condition struct {
	Query string
	Args  []interface{}
}

// String returns the string representation of the condition.
func (c *Condition) String() string {
	argsCount := len(c.Args)
	if argsCount == 1 {
		return fmt.Sprintf(c.Query, "?")
	} else if argsCount > 1 {
		return fmt.Sprintf(c.Query, Placeholders(argsCount))
	} else {
		return c.Query
	}
}
