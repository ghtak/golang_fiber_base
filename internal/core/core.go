package core

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiber(lc fx.Lifecycle, env Env, log *zap.Logger) *fiber.App {
	app := fiber.New()
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info("Server Listen", zap.String("Address", env.ServerAddress))
				go app.Listen(env.ServerAddress)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
	return app
}

func Module() fx.Option {
	return fx.Module(
		"ModuleCore",
		fx.Provide(NewEnv, NewFiber, NewLogger),
		fx.Provide(NewWriteSyncer, NewEncoder, fx.Private),
		fx.Provide(fx.Annotate(NewRouter, fx.ParamTags(`group:"router"`))),
		fx.Invoke(func(*fiber.App) {}),
		fx.Invoke(func(fiber.Router) {}),
	)
}
