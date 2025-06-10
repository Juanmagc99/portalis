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

	service, ok := r.instances[serviceName]
	if !ok {
		return errors.New("service not found")
	}

	inst, ok := service[instanceID]
	if !ok {
		return errors.New("instance not found")
	}

	inst.LastSeen = time.Now()
	service[instanceID] = inst
	return nil
}

func (r *MemRegistry) Deregister(serviceName, instanceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	instances, ok := r.instances[serviceName]
	if !ok {
		return errors.New("service not found")
	}

	if _, ok := instances[instanceID]; !ok {
		return errors.New("instance not found")
	}

	delete(instances, instanceID)

	if len(instances) == 0 {
		delete(r.instances, serviceName)
	}

	return nil
}

func (r *MemRegistry) List(serviceName ...string) ([]model.Instance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []model.Instance

	if len(serviceName) == 1 {
		instances, ok := r.instances[serviceName[0]]
		if !ok {
			return nil, nil
		}
		for _, inst := range instances {
			if time.Since(inst.LastSeen) <= r.ttl {
				result = append(result, inst)
			}
		}
		return result, nil
	}

	for _, services := range r.instances {
		for _, inst := range services {
			if time.Since(inst.LastSeen) <= r.ttl {
				result = append(result, inst)
			}
		}
	}

	return result, nil
}
