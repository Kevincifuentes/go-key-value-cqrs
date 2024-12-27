//go:build debug

package api

import (
	"go-key-value-cqrs/infrastructure/api/config"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type ProfilingServer struct {
	*http.Server
}

func newProfilingServer(serverAddress string) *ProfilingServer {
	return &ProfilingServer{
		&http.Server{
			Addr:    serverAddress,
			Handler: http.DefaultServeMux,
		},
	}
}

func StartProfilingServer(applicationConfig config.Config) {
	serverAddress := applicationConfig.GetDebugServerAddress()
	profilingServer := newProfilingServer(serverAddress)
	go func() {
		log.Printf("Starting Profiling server on %v\n", serverAddress)
		log.Fatal(profilingServer.ListenAndServe())
	}()
}

func init() {
	StartProfilingServer(config.RetrieveConfiguration())
}
