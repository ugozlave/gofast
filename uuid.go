package gofast

import (
	"fmt"
	"sync/atomic"
	"uuid"
)

type UniqueIDGenerator interface {
	Next() string
}

/*
** SequenceIDGenerator
 */

type SequenceIDGenerator struct {
	current atomic.Int64
}

func SequenceIDGeneratorBuilder() Builder[*SequenceIDGenerator] {
	return func(ctx *BuilderContext) *SequenceIDGenerator {
		return &SequenceIDGenerator{}
	}
}

func (g *SequenceIDGenerator) Next() string {
	id := g.current.Add(1)
	return fmt.Sprintf("%d", id)
}

/*
** UUIDGenerator
 */

type UUIDGenerator struct {
}

func UUIDGeneratorBuilder() Builder[*UUIDGenerator] {
	return func(ctx *BuilderContext) *UUIDGenerator {
		return &UUIDGenerator{}
	}
}

func (g *UUIDGenerator) Next() string {
	return uuid.New().String()
}
