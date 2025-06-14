package handler

import (
	"github.com/Juanmagc99/portalis/internal/registry"
)

type RegistryHandler struct {
	Store registry.Registry
}

func New(store registry.Registry) *RegistryHandler {
	return &RegistryHandler{
		Store: store,
	}
}
