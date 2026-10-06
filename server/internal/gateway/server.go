package gateway

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
)

type Config struct {
	ListenAddr string
	BackendAddr string
}

func Start(cfg Config) {

	mux := http.NewServeMux()
	mux.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, cfg.BackendAddr)
	})

	// Live profiling for cloud deployments where localhost pprof is not
	// reachable. Off by default; requests must carry the gateway key.
	if os.Getenv("PPROF_ENABLED") == "1" {
		guard := func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				if !Authorize(r) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				h(w, r)
			}
		}
		mux.HandleFunc("/debug/pprof/", guard(pprof.Index))
		mux.HandleFunc("/debug/pprof/profile", guard(pprof.Profile))
		mux.HandleFunc("/debug/pprof/trace", guard(pprof.Trace))
		mux.HandleFunc("/debug/pprof/symbol", guard(pprof.Symbol))
		mux.HandleFunc("/debug/pprof/cmdline", guard(pprof.Cmdline))
		fmt.Println("pprof enabled on /debug/pprof/ (key required)")
	}

	fmt.Println("Gateway running on", cfg.ListenAddr)

	err := http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		panic(err)
	}
}