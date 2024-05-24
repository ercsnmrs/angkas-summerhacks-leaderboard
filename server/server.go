package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// Server represents application server.
type Server struct {
	*http.Server

	driverService      driverService
	leaderboardService leaderboardService
	authenticator      authenticator
	databaseChecker    databaseChecker
	tracing            tracing
	logger             *slog.Logger

	config  Config
	Version Version
}

// New creates new instance of Server.
func New(
	config Config,
	ds driverService,
	ls leaderboardService,
	authenticator authenticator,
	tracing tracing,
	version Version,
	logger *slog.Logger,
) *Server {
	c := config.setDefaults()

	l := logger.With("pkg", "server")
	l.Info("config",
		"addr", c.Addr,
		"read-timeout", c.ReadTimeout.String(),
		"write-timeout", c.WriteTimeout.String(),
		"shutdown-timeout", c.ShutdownTimeout.String(),
	)

	s := &Server{
		driverService:      ds,
		leaderboardService: ls,
		authenticator:      authenticator,
		tracing:            tracing,
		Version:            version,
		logger:             l,
	}
	s.Server = &http.Server{
		Addr:         c.Addr,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		Handler:      s.Routes(),
	}
	return s
}

// Routes setups middlewares and route endpoints.
func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(
		s.tracing.Middleware(),
		requestIDMiddleware,
		s.loggingMiddleware,
		s.recoveryMiddleware,
	)

	// Public endpoints
	r.Get("/version", GetVersion(s.Version))
	r.Get("/healthcheck", HealthCheck(s.databaseChecker))
	//r.Get("/fighters/{id}", GetFighterByID(s.service))

	// Leaderboard Endpoints
	r.Get("/leaderboard/ranking/{id}", GetDriverRating(s.driverService))
	r.Get("/leaderboard/{scope}", GetLeaderboard(s.leaderboardService))

	// Private endpoints
	r.Route("/", func(r chi.Router) {
		r.Use(authMiddleware(s.authenticator))
		//r.Get("/fighters", ListFighters(s.service))
	})

	r.NotFound(noMatchHandler(http.StatusNotFound))
	r.MethodNotAllowed(noMatchHandler(http.StatusMethodNotAllowed))
	return r
}

// Stop shuts down server gracefully with deadline of shutdownTimeout.
func (s *Server) Stop() error {
	timeout := s.config.ShutdownTimeout
	done := make(chan error, 1)
	go func() {
		ctx := context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		s.logger.Info("shutting down gracefully...")
		done <- s.Shutdown(ctx)
		s.logger.Info("shutdown")
	}()
	return <-done
}

// Run starts serving and listening http server with graceful shutdown.
func (s *Server) Run() error {
	s.logger.Info(fmt.Sprintf("running on %s", s.Addr))
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func noMatchHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := errors.New(http.StatusText(status))
		encodeJSONError(w, e, status)
	}
}

type tracing interface {
	Middleware() func(next http.Handler) http.Handler
}

// default server config values.
const (
	defaultAddr            = ":8000"
	defaultShutdownTimeout = time.Second * 5
)

// Config represents server config.
type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

func (c Config) setDefaults() Config {
	if strings.TrimSpace(c.Addr) == "" {
		c.Addr = defaultAddr
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = defaultShutdownTimeout
	}
	return c
}
