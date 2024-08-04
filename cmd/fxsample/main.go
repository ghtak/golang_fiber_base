package fxsample

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
)

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

type Route interface {
	http.Handler
	Pattern() string
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

func (*EchoHandler) Pattern() string {
	return "/echo"
}

type HelloHandler struct {
	log *zap.Logger
}

func NewHelloHandler(log *zap.Logger) *HelloHandler {
	return &HelloHandler{log}
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error("Failed to read request", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := fmt.Fprintf(w, "Hello, %s\n", body); err != nil {
		h.log.Error("Failed to write response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (*HelloHandler) Pattern() string {
	return "/hello"
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

func main() {
	fx.New(
		fx.Provide(
			NewHttpServer,
			//fx.Annotate(
			//	NewEchoHandler,
			//	fx.As(new(Route)),
			//	fx.ResultTags(`name:"echo"`)),
			//fx.Annotate(
			//	NewHelloHandler,
			//	fx.As(new(Route)),
			//	fx.ResultTags(`name:"hello"`)),
			//fx.Annotate(
			//	NewServeMux,
			//	fx.ParamTags(`name:"echo"`, `name:"hello"`)),
			AsRoute(NewEchoHandler),
			AsRoute(NewHelloHandler),
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`)),
			zap.NewExample,
		),
		fx.Invoke(func(server *http.Server) {}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
