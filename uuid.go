package gofast

import (
	"fmt"
	"sync/atomic"
)

type UniqueIDGenerator interface {
	Next() string
}

/*
** SequenceIDGenerator
 */

type SequenceIDGenerator struct {
	current int64
}

func SequenceIDGeneratorBuilder() Builder[*SequenceIDGenerator] {
	return func(ctx *BuilderContext) *SequenceIDGenerator {
		return &SequenceIDGenerator{}
	}
}

func (g *SequenceIDGenerator) Next() string {
	id := atomic.AddInt64(&g.current, 1)
	return fmt.Sprintf("%d", id)
}
