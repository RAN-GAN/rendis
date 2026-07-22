package rendis

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func New(url, key string) (*Client, error) {
	header := http.Header{}
	header.Set("X-RENDIS-Key", key)

	conn, _, err := websocket.DefaultDialer.Dial(url+"/connect", header)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) send(data []byte) error {
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (c *Client) receive() ([]byte, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return data, nil
}
