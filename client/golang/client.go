package rendis

import (
	"errors"
	"strconv"
)

func (c *Client) Ping() error {
	if err := c.send(encode("PING")); err != nil {
		return err
	}

	resp, err := c.receive()
	if err != nil {
		return err
	}

	value, err := decode(resp)
	if err != nil {
		return err
	}

	pong, ok := value.(string)
	if !ok {
		return ErrInvalidRESP
	}

	if pong != "PONG" {
		return errors.New("unexpected response: " + pong)
	}

	return nil
}
func (c *Client) Set(key, value string) error {
	data := encode("SET", key, value)

	if err := c.send(data); err != nil {
		return err
	}

	resp, err := c.receive()
	if err != nil {
		return err
	}

	_, err = decode(resp)
	return err
}
func (c *Client) Get(key string) (string, error) {
	data := encode("GET", key)

	if err := c.send(data); err != nil {
		return "", err
	}

	resp, err := c.receive()
	if err != nil {
		return "", err
	}

	value, err := decode(resp)
	s, ok := value.(string)
	if !ok {
		return "", ErrInvalidRESP
	}
	return s, err
}
func (c *Client) Del(key string) (int64, error) {
	data := encode("DEL", key)

	if err := c.send(data); err != nil {
		return -1, err
	}

	resp, err := c.receive()

	if err != nil {
		return -1, err
	}

	value, err := decode(resp)
	if err != nil {
		return -1, err
	}

	deleted, ok := value.(int64)
	if !ok {
		return -1, ErrInvalidRESP
	}

	return deleted, nil
}
func (c *Client) TTL(key string) (int64, error) {
	data := encode("TTL", key)

	if err := c.send(data); err != nil {
		return -1, err
	}

	resp, err := c.receive()

	if err != nil {
		return -1, err
	}

	value, err := decode(resp)
	if err != nil {
		return -1, err
	}

	ttl, ok := value.(int64)
	if !ok {
		return -1, ErrInvalidRESP
	}

	return ttl, nil
}
func (c *Client) Expire(key string, seconds int64) (bool, error) {
	data := encode("EXPIRE", key, strconv.FormatInt(seconds, 10))

	if err := c.send(data); err != nil {
		return false, err
	}

	resp, err := c.receive()
	if err != nil {
		return false, err
	}

	value, err := decode(resp)
	if err != nil {
		return false, err
	}

	result, ok := value.(int64)
	if !ok {
		return false, ErrInvalidRESP
	}

	return result == 1, nil
}
func (c *Client) Exists(key string) (bool, error) {
	data := encode("EXISTS", key)

	if err := c.send(data); err != nil {
		return false, err
	}

	resp, err := c.receive()

	if err != nil {
		return false, err
	}

	value, err := decode(resp)
	if err != nil {
		return false, err
	}

	exists, ok := value.(int64)
	if !ok {
		return false, ErrInvalidRESP
	}

	return exists == 1, nil
}

// Do executes a generic command on the Rendis server
func (c *Client) Do(args ...string) (any, error) {
	data := encode(args...)

	if err := c.send(data); err != nil {
		return nil, err
	}

	resp, err := c.receive()
	if err != nil {
		return nil, err
	}

	return decode(resp)
}
