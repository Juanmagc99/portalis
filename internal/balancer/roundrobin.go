package balancer

import (
	"fmt"
	"sync/atomic"

	"github.com/Juanmagc99/portalis/internal/model"
)

type RoundRobin struct {
	counter uint64
}

func (rr *RoundRobin) Next(instances []model.Instance) string {
	n := uint64(len(instances))
	if n == 0 {
		return ""
	}

	idx := atomic.AddUint64(&rr.counter, 1)
	chosen := instances[int((idx-1)%n)]
	return fmt.Sprintf("http://%s:%d", chosen.Host, chosen.Port)
}
