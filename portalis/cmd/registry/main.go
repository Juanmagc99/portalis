package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Juanmagc99/portalis/internal/handler"
	"github.com/Juanmagc99/portalis/internal/registry"
	"github.com/Juanmagc99/portalis/pkg/httperror"
	"github.com/Juanmagc99/portalis/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = httperror.HandleError
	e.Validator = &validation.CustomValidator{Validator: validator.New()}

	ttl := 30 * time.Second
	evictInterval := 10 * time.Second

	s := registry.NewMemRegistry(ttl)
	stopCh := make(chan struct{})
	s.StartEvictor(evictInterval, stopCh)

	handler.NewRoutes(e, s)

	go func() {
		e.Logger.Fatal(e.Start("localhost:8080"))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	close(stopCh)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("error closing server:", err)
	}
}
