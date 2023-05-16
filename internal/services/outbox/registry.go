package outbox

import (
	"errors"
	"fmt"
	"sync"
)

type registry struct {
	mtx         sync.RWMutex
	jobRegistry map[string]Job
}

func newRegistry() *registry {
	return &registry{
		mtx:         sync.RWMutex{},
		jobRegistry: map[string]Job{},
	}
}

func (r *registry) set(job Job) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if _, ok := r.jobRegistry[job.Name()]; ok {
		return errors.New("job with same name already exists")
	}
	r.jobRegistry[job.Name()] = job
	return nil
}

func (r *registry) get(name string) (Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if _, ok := r.jobRegistry[name]; !ok {
		return nil, fmt.Errorf("job with name %q does not exists", name)
	}
	return r.jobRegistry[name], nil
}
