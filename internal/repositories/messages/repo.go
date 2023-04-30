package messagesrepo

import (
	"fmt"

	"github.com/Pickausernaame/chat-service/internal/store"
)

//go:generate options-gen -out-filename=repo_options.gen.go -from-struct=Options
type Options struct {
	db *store.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Options
}

func New(opts Options) (*Repo, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating options: %v", err)
	}
	return &Repo{Options: opts}, nil
}