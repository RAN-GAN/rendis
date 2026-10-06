package gateway

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"slices"
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

	// ws -> tcp: stream each frame through one reused buffer instead of
	// ReadMessage, which allocates a fresh slice per frame.
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			_, r, err := ws.NextReader()
			if err != nil {
				return
			}
			for {
				n, rerr := r.Read(buf)
				if n > 0 {
					if _, err := tcp.Write(buf[:n]); err != nil {
						return
					}
				}
				if rerr != nil {
					break
				}
			}
		}
	}()

	// tcp -> ws: frame one RESP reply per WebSocket message, building it
	// in a reused buffer.
	reader := bufio.NewReader(tcp)
	var out []byte
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			break
		}
		if len(line) <= 2 {
			continue
		}
		out = append(out[:0], line...)

		if line[0] == '$' {
			length, err := strconv.Atoi(string(line[1 : len(line)-2]))
			if err != nil {
				continue
			}
			if length != -1 {
				start := len(out)
				out = slices.Grow(out, length+2)[:start+length+2]
				if _, err := io.ReadFull(reader, out[start:]); err != nil {
					break
				}
			}
		}

		if err := ws.WriteMessage(websocket.BinaryMessage, out); err != nil {
			break
		}
	}
}
