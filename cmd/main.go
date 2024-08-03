package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/app"
	"github.com/golang_fiber_base/internal/log"
	"github.com/golang_fiber_base/internal/module/simple"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

func _main() {
	logConfig := log.Config{
		Level:       "DEBUG",
		RollingFile: "./logs/app.log",
	}
	fiberConfig := fiber.Config{}
	app_ := app.NewApp(logConfig, fiberConfig)

	app_.UseFromLocals()
	app_.Logger.Info("Info Logs")
	app_.Logger.Debug("Debug Logs",
		// Structured context as strongly typed Field values.
		zap.String("string", "stringvalue"),
		zap.Int("int", 3),
		zap.Duration("duration", time.Second),
	)

	simpleModule := app_.Fiber.Group("/simple")
	simple.RegisterRoute(simpleModule)

	err := app_.Listen(":3003")
	app_.Logger.Info(err.Error())
}

func NewHttpServer(lc fx.Lifecycle) *http.Server {
	srv := &http.Server{Addr: ":8009"}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", srv.Addr)
				if err != nil {
					return err
				}
				fmt.Println("Starting HTTP Server at", srv.Addr)
				go srv.Serve(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		})
	return srv
}

func main() {
	fx.New(
		fx.Provide(NewHttpServer),
		fx.Invoke(func(server *http.Server) {}),
	).Run()
}
