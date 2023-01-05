package client

type Option = func(c *config)

func WithProxyUrl(url string) Option {
	return func(c *config) {
		c.proxyUrl = url
	}
}

func WithUseSession(useSession bool) Option {
	return func(c *config) {
		c.useSession = useSession
	}
}

func parseOptions(options []Option) *config {
	c := newConfig()

	for _, o := range options {
		o(c)
	}

	return c
}
