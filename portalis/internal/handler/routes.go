package handler

import (
	"github.com/Juanmagc99/portalis/internal/registry"
	"github.com/labstack/echo/v4"
)

func NewRoutes(e *echo.Echo, store registry.Registry) {
	rh := NewRegistryHandler(store)

	e.POST("/api/registries", rh.Register)
}
