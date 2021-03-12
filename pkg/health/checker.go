package health

import (
	"sync"
)

// Checker is a interface used to provide an indication of application health.
type Checker interface {
	Name() string
	Check() Health
}

// CompositeChecker aggregate a list of Checkers
type CompositeChecker struct {
	checkers []Checker
	info     map[string]interface{}
}

// NewCompositeChecker creates a new CompositeChecker
func NewCompositeChecker() CompositeChecker {
	return CompositeChecker{}
}

// AddInfo adds a info value to the Info map
func (c *CompositeChecker) AddInfo(key string, value interface{}) *CompositeChecker {
	if c.info == nil {
		c.info = make(map[string]interface{})
	}

	c.info[key] = value

	return c
}

// AddChecker add a Checker to the aggregator
func (c *CompositeChecker) AddChecker(checker Checker) {
	c.checkers = append(c.checkers, checker)
}

// Check returns the combination of all checkers added
// if some check is not up, the combined is marked as down
func (c CompositeChecker) Check() Health {
	health := NewHealth()
	health.Up()

	healths := make(map[string]interface{})

	type state struct {
		h    Health
		name string
	}
	ch := make(chan state, len(c.checkers))
	var wg sync.WaitGroup
	for _, c := range c.checkers {
		wg.Add(1)
		go func(name string, f func() Health) {
			ch <- state{h: f(), name: name}
			wg.Done()
		}(c.Name(), c.Check)
	}
	wg.Wait()
	close(ch)

	for s := range ch {
		if !s.h.IsUp() && !health.IsDown() {
			health.Down()
		}
		healths[s.name] = s.h
	}

	health.info = healths

	// Extra Info
	for key, value := range c.info {
		health.AddInfo(key, value)
	}
	return health
}
