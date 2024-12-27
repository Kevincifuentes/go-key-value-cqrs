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

func newProfilingServer(applicationConfig config.Config) *ProfilingServer {
	return &ProfilingServer{
		&http.Server{
			Addr:    applicationConfig.GetDebugServerAddress(),
			Handler: http.DefaultServeMux,
		},
	}
}

func StartProfilingServer(applicationConfig config.Config) {
	profilingServer := newProfilingServer(applicationConfig)
	go func() {
		log.Printf("Starting Profiling server on %v\n", applicationConfig.DebugServerPort)
		log.Fatal(profilingServer.ListenAndServe())
	}()
}

func init() {
	StartProfilingServer(config.RetrieveConfiguration())
}
