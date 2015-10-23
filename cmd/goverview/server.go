package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vbatoufflet/goverview/pkg/logger"
	"github.com/vbatoufflet/goverview/pkg/poller"
)

// Server represents a server struct
type Server struct {
	config Config
	poller *poller.Poller
	router *gin.Engine
}

// NewServer returns a new server instance
func NewServer(c Config) *Server {
	return &Server{
		config: c,
	}
}

// Run starts the server processing
func (s *Server) Run() error {
	var err error

	logger.Info("server", "service started")

	// Initialize poller
	s.poller = poller.NewPoller(s.config.Poller)

	// Initialize HTTP router
	gin.SetMode(gin.ReleaseMode)

	s.router = gin.New()
	s.router.RedirectTrailingSlash = true

	// Set static files handlers
	s.router.Any("/", s.handleStatic)
	s.router.Any("/static/*path", s.handleStatic)

	// Set API handlers
	s.router.Any("/api/hosts/:name", s.handleAPIHost)
	s.router.Any("/api/hosts/", s.handleAPIHosts)
	s.router.Any("/api/groups/", s.handleAPIGroups)
	s.router.Any("/api/nodes/", s.handleAPINodes)
	s.router.Any("/api/search/", s.handleAPISearch)

	// Start HTTP router
	logger.Info("server", "listening on %q", s.config.BindAddr)

	err = s.router.Run(s.config.BindAddr)
	if err != nil {
		logger.Error("server", "failed to initialize router: %s", err)
		return err
	}

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	s.poller.Shutdown()

	logger.Info("server", "service stopped")

	os.Exit(0)
}

// Refresh triggers a poller refresh
func (s *Server) Refresh() {
	s.poller.Refresh()
}
