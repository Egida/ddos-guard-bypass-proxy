package worker

import "time"

type Option = func(c *config)

func WithFlareSolverrUrl(url string) Option {
	return func(c *config) {
		c.flareSolverrUrl = url
	}
}

func WithQueryTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.queryTimeout = timeout
	}
}

func parseOptions(options []Option) *config {
	c := newConfig()

	for _, o := range options {
		o(c)
	}

	return c
}
