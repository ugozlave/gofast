package gofast

import (
	"context"

	"github.com/ugozlave/cargo"
)

type ContextKey string

const (
	CtxName      ContextKey = "Name"
	CtxRequestId ContextKey = "RequestId"
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

func (c *BuilderContext) RequestID() string {
	v, ok := c.Value(CtxRequestId).(string)
	if !ok {
		return ""
	}
	return v
}

func (c *BuilderContext) Get[T any](lt Lifetime) T {
	return Get[T](c, lt)
}

func (c *BuilderContext) MustGet[T any](lt Lifetime) T {
	return MustGet[T](c, lt)
}

func (c *BuilderContext) GetLogger[S any](lt Lifetime) Logger {
	return GetLogger[S](c, lt)
}

func (c *BuilderContext) MustGetLogger[S any](lt Lifetime) Logger {
	return MustGetLogger[S](c, lt)
}

func (c *BuilderContext) GetConfig[C any](lt Lifetime) Config[C] {
	return GetConfig[C](c, lt)
}

func (c *BuilderContext) MustGetConfig[C any](lt Lifetime) Config[C] {
	return MustGetConfig[C](c, lt)
}

func (c *BuilderContext) All[T any](lt Lifetime) []T {
	return All[T](c, lt)
}
