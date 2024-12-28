package web

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.quinn.io/ccf/assets"
	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/internal/router"
)

//go:embed public
var assetsFS embed.FS

type Server struct {
	e *echo.Echo
	b *boot.Boot
}

func NewServer(b *boot.Boot) *Server {
	e := echo.New()
	e.Use(echomiddleware.Logger())

	// Register routes from generated code
	router.RegisterRoutes(e)

	// Attach public assets
	assets.Attach(
		e,
		"public",
		"internal/web/public",
		assetsFS,
		os.Getenv("USE_EMBEDDED_ASSETS") == "true",
	)

	e.Debug = true
	e.Use(middleware.BootContext(b))
	e.Use(middleware.EchoContext(e))
	e.HTTPErrorHandler = middleware.HTTPErrorHandler
	addRoutes(e)

	return &Server{e: e, b: b}
}

func (s *Server) Run(ctx context.Context) error {
	// Start server in goroutine
	go func() {
		fmt.Println("Server starting on http://localhost:3000")
		if err := s.e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
