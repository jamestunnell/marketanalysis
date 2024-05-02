package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

const (
	ShutdownTimeout = 15 * time.Second
)

func New(port int) *Server {
	httpSrv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		WriteTimeout:      time.Minute,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: time.Minute,
		IdleTimeout:       time.Minute,
	}
	r := mux.NewRouter()

	httpSrv.Handler = r

	srv := &Server{
		httpServer: httpSrv,
		router:     r,
	}

	return srv
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) Start() {
	go s.listen()
}

func (s *Server) Stop() {
	ctx, _ := context.WithTimeout(context.Background(), ShutdownTimeout)

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("http server: failed to shut down")
	}
}

func (s *Server) listen() {
	log.Info().
		Str("addr", s.httpServer.Addr).
		Msg("http server: started")

	err := s.httpServer.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http server: listen error")
	}

	log.Info().Msg("http server: stopped")
}
