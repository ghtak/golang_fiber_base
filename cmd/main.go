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
	"io"
	"net"
	"net/http"
	"os"
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

func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8009", Handler: mux}
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

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (*EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to handler request:", err)
	}
}

func NewServeMux(echo *EchoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/echo", echo)
	return mux
}

func main() {
	fx.New(
		fx.Provide(
			NewHttpServer,
			NewEchoHandler,
			NewServeMux,
		),
		fx.Invoke(func(server *http.Server) {}),
	).Run()
}
