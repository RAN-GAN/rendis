package rendis

import (
	"bytes"
	"errors"
	"strconv"
)

var ErrInvalidRESP = errors.New("invalid RESP message")

func encode(args ...string) []byte {
	var buf bytes.Buffer
	suffix := "\r\n"
	encoded := "*" + strconv.Itoa(len(args)) + suffix
	buf.WriteString(encoded)
	for _, arg := range args {
		buf.WriteString(bulkString(arg))
	}
	return buf.Bytes()
}

func decode(message []byte) (any, error) {
	if len(message) == 0 {
		return nil, ErrInvalidRESP
	}
	prefix := message[0]
	switch prefix {
	case '+':
		value := string(message[1 : len(message)-2])
		return value, nil
	case '$':
		length, err := strconv.Atoi(string(message[1:bytes.IndexByte(message, '\r')]))
		if err != nil {
			return nil, ErrInvalidRESP
		}
		if length == -1 {
			return nil, nil
		}
		headerEnd := bytes.Index(message, []byte("\r\n"))
		if headerEnd == -1 {
			return nil, ErrInvalidRESP
		}
		start := headerEnd + 2
		end := start + length

		if end+2 != len(message) {
			return nil, ErrInvalidRESP
		}

		value := string(message[start:end])
		return value, nil

	case '-':
		errorMsg := string(message[1 : len(message)-2])
		return nil, errors.New(errorMsg)

	case ':':

		value, err := strconv.ParseInt(string(message[1:len(message)-2]), 10, 64)
		if err != nil {
			return nil, err
		}
		return value, nil

	default:
		return nil, ErrInvalidRESP
	}
}

func bulkString(value string) string {
	return "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"
}
