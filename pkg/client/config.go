package client

type config struct {
	proxyUrl   string
	useSession bool
}

func newConfig() *config {
	return &config{
		proxyUrl:   DefaultProxyUrl,
		useSession: DefaultUseSession,
	}
}
