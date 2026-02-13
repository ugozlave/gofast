package gofast

import (
	"context"

	"github.com/ugozlave/cargo"
)

type ContextKey string

const (
	CtxName        ContextKey = "Name"
	CtxEnvironment ContextKey = "Environment"
	CtxRequestId   ContextKey = "RequestId"
)

type BuilderContext struct {
	context.Context
	container *cargo.Container
}

func NewBuilderContext(ctx context.Context, container *cargo.Container) *BuilderContext {
	return &BuilderContext{
		Context:   ctx,
		container: container,
	}
}

func (c *BuilderContext) Name() string {
	v, ok := c.Value(CtxName).(string)
	if !ok {
		return ""
	}
	return v
}

func (c *BuilderContext) Environment() string {
	v, ok := c.Value(CtxEnvironment).(string)
	if !ok {
		return ""
	}
	return v
}

func (c *BuilderContext) RequestID() string {
	v, ok := c.Value(CtxRequestId).(string)
	if !ok {
		return ""
	}
	return v
}
