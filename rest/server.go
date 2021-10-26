package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jrpalma/pwdhash/config"
	"github.com/jrpalma/pwdhash/logs"
)

func NewServer(conf config.Config, log logs.Logger) *Server {
	server := &Server{log: log}
	server.handler = newHandler(conf, log, server.Shutdown)

	server.mux = http.NewServeMux()
	server.mux.HandleFunc(v1+"/hash", server.handler.newHash)
	server.mux.HandleFunc(v1+"/hash/", server.handler.checkHash)
	server.mux.HandleFunc(v1+"/stats", server.handler.stats)
	server.mux.HandleFunc(v1+"/shutdown", server.handler.shutdown)

	server.httpServer = &http.Server{
		Addr:    conf.ServerAddress,
		Handler: server,
	}

	return server
}

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	log        logs.Logger
	handler    *handler
	requestID  int
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rid := r.Header.Get("X-Request-ID")
	if rid == "" {
		rid = strconv.Itoa(s.requestID)
		r.Header.Set("X-Request-ID", rid)
		s.requestID++
	}

	ctx := context.WithValue(r.Context(), ridKey, rid)

	start := time.Now()
	s.mux.ServeHTTP(w, r.WithContext(ctx))
	duration := time.Since(start)

	s.log.Debugf("%s(%v) %s %s %v", ridKey, rid, r.Method, r.URL.Path, duration)
}

func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		s.log.Errorf("Failed to start server: %v", err)
	}
	return err
}

func (s *Server) Shutdown() error {
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		s.log.Errorf("Failed to shutdown server: %v", err)
	}
	return err
}

const (
	ridKey = "REQID"
	v1     = "/api/v1"
)
