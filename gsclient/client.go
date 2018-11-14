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

func (c *Client) ValidateToken() error {
	r := Request{
		uri:    "/validate_token/",
		method: "GET",
	}

	return r.execute(*c, nil)
}
