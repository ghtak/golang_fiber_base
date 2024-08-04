package application

import (
	"github.com/golang_fiber_base/internal/application/hello"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"ModuleApplication",
		hello.Module(),
	)
}
