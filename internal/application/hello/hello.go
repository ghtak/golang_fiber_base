package hello

import (
	"github.com/golang_fiber_base/internal/core"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"ModuleHello",
		fx.Provide(core.AsRouter(NewHelloController)),
	)
}
