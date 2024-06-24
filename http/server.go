package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MBvisti/mortenvistisen/http/router"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/gorilla/csrf"
)

type Server struct {
	router *router.Router
	host   string
	port   string
	cfg    config.Cfg
	srv    *http.Server
}

func NewServer(
	router *router.Router,
	logger *slog.Logger,
	cfg config.Cfg,
) Server {
	host := cfg.App.ServerHost
	port := cfg.App.ServerPort
	isProduction := cfg.App.Environment == "production"

	srv := &http.Server{
		Addr: fmt.Sprintf("%v:%v", host, port),
		Handler: csrf.Protect(
			[]byte(cfg.Auth.CsrfToken), csrf.Secure(isProduction))(router.GetInstance()),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return Server{
		router,
		host,
		port,
		cfg,
		srv,
	}
}

func (s *Server) Start() {
	slog.Error("starting server on", "host", s.host, "port", s.port)

	// Start server
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()

	toCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Print("initiating shutdown")
	err := s.srv.Shutdown(toCtx)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("shutdown complete")
}
