package gateway

import (
	"net/http"
	"os"
	"strings"
)

func Authorize(r *http.Request) bool {
	key := os.Getenv("KEY")

	if r.Header.Get("x-rendis-key") == key {
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
