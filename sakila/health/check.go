package health

import "github.com/InVisionApp/go-health"

// Check is a health check.
type Check struct {
	Name    string
	Checker health.ICheckable
}

// Checks are the health checks.
type Checks struct {
	DB    *Check
	Cache *Check
}
