package handler

import (
	"github.com/Juanmagc99/portalis/internal/registry"
	"github.com/labstack/echo/v4"
)

func NewRoutes(e *echo.Echo, store registry.Registry) {
	rh := NewRegistryHandler(store)

	a := e.Group("/api")

	a.POST("/register", rh.Register)
	a.PUT("/heartbeat", rh.Heartbeat)
	a.DELETE("/deregister", rh.Deregister)
	a.GET("/services", rh.List)
}
