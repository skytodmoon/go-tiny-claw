package engine

import (
	"log"

	"github.com/skytodmoon/go-tiny-claw/internal/context"
	"github.com/skytodmoon/go-tiny-claw/internal/provider"
	"github.com/skytodmoon/go-tiny-claw/internal/tools"
)

type Engine struct {
	provider   provider.Provider
	registry   *tools.Registry
	ctxManager *context.ContextManager
}

func NewEngine() *Engine {
	return &Engine{}
}

func NewAgentEngine(p provider.Provider, r *tools.Registry, c *context.ContextManager) *Engine {
	return &Engine{
		provider:   p,
		registry:   r,
		ctxManager: c,
	}
}

func (e *Engine) Run(task string) error {
	log.Printf("Engine is running task: %s", task)
	return nil
}
