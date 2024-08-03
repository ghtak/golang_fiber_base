package main

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/golang_fiber_base/internal/app"
	"github.com/golang_fiber_base/internal/log"
	"github.com/golang_fiber_base/internal/module/simple"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"io"
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

func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
	srv := &http.Server{Addr: ":8009", Handler: mux}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", srv.Addr)
				if err != nil {
					return err
				}
				log.Info("Starting HTTP Server at", zap.String("addr", srv.Addr))
				go srv.Serve(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return srv.Shutdown(ctx)
			},
		})
	return srv
}

type EchoHandler struct {
	log *zap.Logger
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log}
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		h.log.Warn("Failed to handler request:", zap.Error(err))
	}
}

type Route interface {
	http.Handler
	Pattern() string
}

func (*EchoHandler) Pattern() string {
	return "/echo"
}

func NewServeMux(route Route) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(route.Pattern(), route)
	return mux
}

func main() {
	fx.New(
		fx.Provide(
			NewHttpServer,
			fx.Annotate(
				NewEchoHandler,
				fx.As(new(Route))),
			NewServeMux,
			zap.NewExample,
		),
		fx.Invoke(func(server *http.Server) {}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
