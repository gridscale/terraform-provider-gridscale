package gsclient

type Client struct {
	cfg *Config
}

func NewClient(c *Config) *Client {
	client := &Client{
		cfg: c,
	}

	return client
}
