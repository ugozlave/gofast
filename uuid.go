package gofast

import "fmt"

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
	g.current++
	return fmt.Sprintf("%d", g.current)
}
