package handler

import (
	"fmt"
	"net/http"

	"github.com/Juanmagc99/portalis/internal/model"
	"github.com/Juanmagc99/portalis/internal/registry"
	"github.com/labstack/echo/v4"
)

type RegistryHandler struct {
	Store registry.Registry
}

func NewRegistryHandler(store registry.Registry) *RegistryHandler {
	return &RegistryHandler{
		Store: store,
	}
}

func (h *RegistryHandler) Register(c echo.Context) error {
	var req model.Instance
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.Store.Register(req); err != nil {
		return err
	}

	m := fmt.Sprintf("Added instance: %+v", req)
	return c.JSON(http.StatusOK, echo.Map{
		"message": m,
	})
}
