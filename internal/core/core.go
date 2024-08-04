package core

import (
	"context"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewFiber(
	lc fx.Lifecycle,
	env Env,
	log *zap.Logger,
	errorHandler ErrorHandler,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "fiber_base",
		ErrorHandler: errorHandler,
	})
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log,
	}))
	//app.Use(swagger.New())
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info("Server Listen", zap.String("Address", env.ServerAddress))
				go func() {
					err := app.Listen(env.ServerAddress)
					if err != nil {
						log.Error("Listen Error", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
	return app
}

type Param struct {
	fx.In

	App      *fiber.App
	Logger   *zap.Logger
	Env      Env
	Database *sqlx.DB
}

func Module() fx.Option {
	return fx.Module(
		"ModuleCore",
		fx.Provide(NewEnv, NewErrorHandler, NewFiber, NewLogger, NewDatabase),
		fx.Provide(NewWriteSyncer, NewEncoder, fx.Private),
		fx.Provide(fx.Annotate(NewRouter, fx.ParamTags(`group:"router"`))),
		fx.Invoke(func(*fiber.App) {}),
		fx.Invoke(func(fiber.Router) {}),
	)
}
