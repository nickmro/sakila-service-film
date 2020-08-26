package health

import (
	"time"

	"github.com/InVisionApp/go-health/v2"
)

// Checker is a service health checker.
type Checker struct {
	*health.Health
}

const intervalDuration = time.Second * 5

// NewChecker returns a new health checker.
func NewChecker(c *Checks) (*Checker, error) {
	checker := health.New()
	checker.DisableLogging()

	if err := checker.AddCheck(&health.Config{
		Name:     c.DB.Name,
		Checker:  c.DB.Checker,
		Interval: intervalDuration,
		Fatal:    true,
	}); err != nil {
		return nil, err
	}

	if err := checker.AddCheck(&health.Config{
		Name:     c.Cache.Name,
		Checker:  c.Cache.Checker,
		Interval: intervalDuration,
		Fatal:    true,
	}); err != nil {
		return nil, err
	}

	return &Checker{Health: checker}, nil
}
