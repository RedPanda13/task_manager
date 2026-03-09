package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	Register(engine *gin.Engine)
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg Config, registrars ...RouteRegistrar) *Server {
	engine := gin.New()

	for _, registrar := range registrars {
		registrar.Register(engine)
	}

	return &Server{
		httpServer: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.HTTP.Port),
			Handler:           engine,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
