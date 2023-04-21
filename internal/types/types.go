package types

import (
	"github.com/google/uuid"
)

type ResourceTyper interface {
	~[16]byte
}

func Parse[T ResourceTyper](s string) (T, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return T(uuid.Nil), err
	}
	return T(u), nil
}

func MustParse[T ResourceTyper](s string) T {
	res, err := Parse[T](s)
	if err != nil {
		panic(err)
	}
	return res
}
