package gateway

import (
	"fmt"
	"net/http"
)

type Config struct {
	ListenAddr string
	BackendAddr string
}

func Start(cfg Config) {

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, cfg.BackendAddr)
	})

	fmt.Println("Gateway running on", cfg.ListenAddr)

	err := http.ListenAndServe(cfg.ListenAddr, nil)
	if err != nil {
		panic(err)
	}
}