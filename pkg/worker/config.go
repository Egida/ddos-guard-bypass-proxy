package worker

import (
	"os"
	"time"
)

type config struct {
	flareSolverrUrl string
	queryTimeout    time.Duration
}

func newConfig() *config {
	return &config{
		flareSolverrUrl: flareSolverrUrl(),
		queryTimeout:    DefaultQueryTimeout,
	}
}

func flareSolverrUrl() string {
	flareSolverrUrl := os.Getenv("FLARESOLVERR_URL")
	if flareSolverrUrl == "" {
		return DefaultFlareSolverrUrl
	}

	return flareSolverrUrl
}
