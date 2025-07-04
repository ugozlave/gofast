package gofast

import (
	"fmt"
	"sync/atomic"
)

type UniqueIDGenerator interface {
	Next() string
}

type SequenceIDGenerator struct {
	current int64
}

func NewSequenceIDGenerator() *SequenceIDGenerator {
	return &SequenceIDGenerator{}
}

func (g *SequenceIDGenerator) Next() string {
	id := atomic.AddInt64(&g.current, 1)
	return fmt.Sprintf("%d", id)
}
