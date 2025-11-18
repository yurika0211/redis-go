package core

import (
	"net"
)

type Client struct {
	Addr string
	Conn net.Conn
}

/**
 * NewClient creates a new Client instance.
 */
func NewClient(addr string) *Client {
	return &Client{
		Addr: addr,
	}
}

/**
 * 启动客户端
 * @param c *Client 客户端
 */
func StartClient(c *Client) {
	if err := c.Connect(); err != nil {
		panic(err)
	}
}

/**
 * Connect establishes a connection to the server.
 */
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

/**
 * Close the client connection
 * @param c *Client 客户端
 */
func CloseClient(c *Client) {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
