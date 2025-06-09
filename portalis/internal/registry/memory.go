package registry

import (
	"errors"
	"sync"
	"time"

	"github.com/Juanmagc99/portalis/internal/model"
)

type MemRegistry struct {
	mu        sync.Mutex
	instances map[string]map[string]model.Instance
	ttl       time.Duration
}

func NewMemRegistry(ttl time.Duration) *MemRegistry {
	return &MemRegistry{
		instances: make(map[string]map[string]model.Instance),
		ttl:       ttl,
	}
}

func (r *MemRegistry) Register(inst model.Instance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.instances[inst.ServiceName]; !ok {
		r.instances[inst.ServiceName] = make(map[string]model.Instance)
	}

	inst.LastSeen = time.Now()
	r.instances[inst.ServiceName][inst.InstanceID] = inst
	return nil
}

func (r *MemRegistry) Heartbeat(serviceName, instanceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	services, ok := r.instances[serviceName]
	if !ok {
		return errors.New("service not found")
	}

	inst, ok := services[instanceID]
	if !ok {
		return errors.New("instance not found")
	}

	inst.LastSeen = time.Now()
	services[instanceID] = inst
	return nil
}
