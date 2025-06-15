package registry

import (
	"errors"
	"sync"
	"time"

	"github.com/Juanmagc99/portalis/internal/model"
)

type MemRegistry struct {
	rm       sync.Mutex
	services map[string]map[string]model.Instance
	ttl      time.Duration
}

func NewMemRegistry(ttl time.Duration) *MemRegistry {
	return &MemRegistry{
		services: make(map[string]map[string]model.Instance),
		ttl:      ttl,
	}
}

func (r *MemRegistry) StartEvictor(evictInterval time.Duration, stopCh <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(evictInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.evictExpired()
			case <-stopCh:
				return
			}
		}
	}()
}

func (r *MemRegistry) evictExpired() {
	rm := &r.rm
	rm.Lock()
	defer rm.Unlock()

	now := time.Now()
	for serviceName, instances := range r.services {
		for id, inst := range instances {
			if now.Sub(inst.LastSeen) > r.ttl {
				delete(instances, id)
			}
		}
		if len(instances) == 0 {
			delete(r.services, serviceName)
		}
	}
}

func (r *MemRegistry) Register(inst model.Instance) error {
	r.rm.Lock()
	defer r.rm.Unlock()

	if _, exists := r.services[inst.ServiceName]; !exists {
		r.services[inst.ServiceName] = make(map[string]model.Instance)
	}
	inst.LastSeen = time.Now()
	r.services[inst.ServiceName][inst.InstanceID] = inst
	return nil
}

func (r *MemRegistry) Heartbeat(serviceName, instanceID string) error {
	r.rm.Lock()
	defer r.rm.Unlock()

	instances, exists := r.services[serviceName]
	if !exists {
		return errors.New("service not found")
	}

	inst, exists := instances[instanceID]
	if !exists {
		return errors.New("instance not found")
	}
	inst.LastSeen = time.Now()
	instances[instanceID] = inst
	return nil
}

func (r *MemRegistry) Deregister(serviceName, instanceID string) error {
	r.rm.Lock()
	defer r.rm.Unlock()

	instances, exists := r.services[serviceName]
	if !exists {
		return errors.New("service not found")
	}
	if _, exists = instances[instanceID]; !exists {
		return errors.New("instance not found")
	}
	delete(instances, instanceID)
	if len(instances) == 0 {
		delete(r.services, serviceName)
	}
	return nil
}

func (r *MemRegistry) List(serviceNames ...string) ([]model.Instance, error) {
	r.rm.Lock()
	defer r.rm.Unlock()

	now := time.Now()
	var result []model.Instance

	if len(serviceNames) > 0 {
		for _, svc := range serviceNames {
			if instances, exists := r.services[svc]; exists {
				for _, inst := range instances {
					if now.Sub(inst.LastSeen) <= r.ttl {
						result = append(result, inst)
					}
				}
			}
		}
		return result, nil
	}

	for _, instances := range r.services {
		for _, inst := range instances {
			if now.Sub(inst.LastSeen) <= r.ttl {
				result = append(result, inst)
			}
		}
	}
	return result, nil
}
