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
			fmt.Println(data)
			if len(data) > 0 {
				_, err = tcp.Write(data)
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
		rawMsg := []byte(line)

		lineStr := strings.TrimSuffix(line, "\r\n")
		if len(lineStr) == 0 {
			continue
		}

		prefix := lineStr[0]
		content := lineStr[1:]

		if prefix == '$' {
			length, err := strconv.Atoi(content)
			if err != nil {
				continue
			}
			if length != -1 {
				buf := make([]byte, length)
				if _, err := io.ReadFull(reader, buf); err != nil {
					break
				}
				crlf := make([]byte, 2)
				io.ReadFull(reader, crlf)

				rawMsg = append(rawMsg, buf...)
				rawMsg = append(rawMsg, crlf...)
			}
		}

		err = ws.WriteMessage(websocket.BinaryMessage, rawMsg)
		if err != nil {
			break
		}
	}
}
