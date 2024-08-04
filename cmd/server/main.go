package main

import (
	"github.com/golang_fiber_base/internal/application"
	"github.com/golang_fiber_base/internal/core"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		core.Module(),
		application.Module(),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
