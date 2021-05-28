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
func NewChecker(checks []*Check) (*Checker, error) {
	checker := health.New()
	checker.DisableLogging()

	for i := range checks {
		if err := checker.AddCheck(&health.Config{
			Name:     checks[i].Name,
			Checker:  checks[i].Checker,
			Interval: intervalDuration,
			Fatal:    true,
		}); err != nil {
			return nil, err
		}
	}

	return &Checker{Health: checker}, nil
}
