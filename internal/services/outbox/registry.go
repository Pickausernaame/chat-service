package outbox

import (
	"errors"
	"fmt"
)

type registry struct {
	jobRegistry map[string]Job
}

func newRegistry() *registry {
	return &registry{
		jobRegistry: map[string]Job{},
	}
}

func (r *registry) set(job Job) error {
	if _, ok := r.jobRegistry[job.Name()]; ok {
		return errors.New("job with same name already exists")
	}
	r.jobRegistry[job.Name()] = job
	return nil
}

func (r *registry) get(name string) (Job, error) {
	if _, ok := r.jobRegistry[name]; !ok {
		return nil, fmt.Errorf("job with name %q does not exists", name)
	}
	return r.jobRegistry[name], nil
}
