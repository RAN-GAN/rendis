package gateway

import (
	"net/http"
	"os"
	"strings"
)

func Authorize(r *http.Request) bool {
	key := os.Getenv("KEY")

	clientKey := r.Header.Get("x-rendis-key")
	if clientKey == "" {
		clientKey = r.URL.Query().Get("key")
	}
	if clientKey == "" {
		protocols := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ",")
		for i, p := range protocols {
			if strings.TrimSpace(p) == "x-rendis-key" && i+1 < len(protocols) {
				clientKey = strings.TrimSpace(protocols[i+1])
				break
			}
		}
	}

	if clientKey == key {
		return true
	}
	return false
}

func VerifyOrigin(r *http.Request) bool {
	var allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	origin := r.Header.Get("Origin")

	for _, allowed := range allowedOrigins {
		incomingOrigin := strings.TrimSpace(allowed)
		if incomingOrigin == "*" || incomingOrigin == origin {
			return true
		}
	}

	return false
}
