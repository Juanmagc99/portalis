package registry

import "github.com/Juanmagc99/portalis/internal/model"

type Registry interface {
	Register(inst model.Instance) error
	Heartbeat(serviceName, instanceID string) error
	Deregister(serviceName, instanceID string) error
	List(serviceName ...string) ([]model.Instance, error)
}
