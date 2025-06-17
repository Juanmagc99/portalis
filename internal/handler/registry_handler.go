package handler

import (
	"fmt"
	"net/http"

	"github.com/Juanmagc99/portalis/internal/model"
	"github.com/Juanmagc99/portalis/internal/registry"
	"github.com/labstack/echo/v4"
)

type SvcInstRequest struct {
	ServiceName string `json:"serviceName" validate:"required"`
	InstanceID  string `json:"instanceID" validate:"required"`
}

type RegistryHandler struct {
	Store registry.Registry
}

func NewRegistryHandler(store registry.Registry) *RegistryHandler {
	return &RegistryHandler{Store: store}
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
		he := echo.NewHTTPError(http.StatusInternalServerError, "failed to process registration")
		he.Internal = err
		return he
	}

	message := fmt.Sprintf("Added instance: %%+v", req)
	return c.JSON(http.StatusCreated, echo.Map{"message": message})
}

func (h *RegistryHandler) Heartbeat(c echo.Context) error {
	var req SvcInstRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.Store.Heartbeat(req.ServiceName, req.InstanceID)
	if err != nil {
		he := echo.NewHTTPError(http.StatusInternalServerError, "failed to process heartbeat")
		he.Internal = err
		return he
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RegistryHandler) Deregister(c echo.Context) error {
	var req SvcInstRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.Store.Deregister(req.ServiceName, req.InstanceID)
	if err != nil {
		he := echo.NewHTTPError(http.StatusInternalServerError, "failed to process deregister")
		he.Internal = err
		return he
	}

	message := fmt.Sprintf("Deregistered instance: %s from service %s", req.InstanceID, req.ServiceName)
	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (h *RegistryHandler) List(c echo.Context) error {
	svcNames := c.QueryParams()["svc"]
	var instances []model.Instance
	var err error

	switch len(svcNames) {
	case 0:
		instances, err = h.Store.List()
	default:
		instances, err = h.Store.List(svcNames...)
	}

	if err != nil {
		he := echo.NewHTTPError(http.StatusInternalServerError, "failed to list services")
		he.Internal = err
		return he
	}

	return c.JSON(http.StatusOK, instances)
}
