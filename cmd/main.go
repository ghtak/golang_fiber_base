package main

import (
	"github.com/golang_fiber_base/internal/app"
	"github.com/golang_fiber_base/internal/config"
	"github.com/golang_fiber_base/internal/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		config.Module(),
		log.Module(),
		app.Module(),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
