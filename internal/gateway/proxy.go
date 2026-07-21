package gateway

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return VerifyOrigin(r)
	},
}

func translateTextToRESP(text string) []byte {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	args := strings.Fields(text)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return buf.Bytes()
}

func handleConnection(w http.ResponseWriter, r *http.Request, backend string) {
	ok := Authorize(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	tcp, err := net.Dial("tcp", backend)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("-ERR backend unavailable"))
		return
	}
	defer tcp.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				return
			}

			respData := translateTextToRESP(string(data))
			if len(respData) > 0 {
				_, err = tcp.Write(respData)
				if err != nil {
					return
				}
			}
		}
	}()

	reader := bufio.NewReader(tcp)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\r\n")

		if len(line) == 0 {
			continue
		}

		var output string
		prefix := line[0]
		content := line[1:]

		switch prefix {
		case '+':
			output = content
		case '-':
			output = "(error) " + content
		case ':':
			output = "(integer) " + content
		case '$':
			length, err := strconv.Atoi(content)
			if err != nil {
				continue
			}
			if length == -1 {
				output = "(nil)"
			} else {
				buf := make([]byte, length)
				if _, err := io.ReadFull(reader, buf); err != nil {
					break
				}
				crlf := make([]byte, 2)
				io.ReadFull(reader, crlf)

				output = string(buf)
			}
		default:
			output = "unknown response: " + line
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte(output))
		if err != nil {
			break
		}
	}
}
