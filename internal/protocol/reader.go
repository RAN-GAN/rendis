package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSuffix(line, "\r\n")

	countStr, ok := strings.CutPrefix(line, "*")
	if !ok {
		return nil, fmt.Errorf("protocol error: expected array")
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, fmt.Errorf("protocol error: invalid array length")
	}
	if count < 0 {
		return nil, fmt.Errorf("protocol error: invalid array length")
	}

	parts := make([]string, 0, count)

	for i := 0; i < count; i++ {
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSuffix(line, "\r\n")

		lengthStr, ok := strings.CutPrefix(line, "$")
		if !ok {
			return nil, fmt.Errorf("protocol error: expected bulk string")
		}

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("protocol error: invalid bulk string length")
		}
		if length < 0 {
			return nil, fmt.Errorf("protocol error: invalid bulk string length")
		}

		buf := make([]byte, length)

		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}

		crlf := make([]byte, 2)

		_, err = io.ReadFull(reader, crlf)
		if err != nil {
			return nil, err
		}

		if crlf[0] != '\r' || crlf[1] != '\n' {
			return nil, fmt.Errorf("protocol error: expected CRLF")
		}

		parts = append(parts, string(buf))
	}

	return parts, nil
}
