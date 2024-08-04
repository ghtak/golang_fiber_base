package core

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type Route struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

type Router interface {
	Routes() []Route
	Group() fiber.Router
}

func AsRouter(r interface{}) interface{} {
	return fx.Annotate(
		r,
		fx.As(new(Router)),
		fx.ResultTags(`group:"router"`))
}

type MapRouteHandler = func(string, ...fiber.Handler) fiber.Router

func NewRouter(routers []Router) fiber.Router {
	for _, router := range routers {
		group := router.Group()
		m := map[string]MapRouteHandler{
			fiber.MethodGet:     group.Get,
			fiber.MethodHead:    group.Head,
			fiber.MethodPost:    group.Post,
			fiber.MethodPut:     group.Put,
			fiber.MethodPatch:   group.Patch,
			fiber.MethodDelete:  group.Delete,
			fiber.MethodConnect: group.Connect,
			fiber.MethodOptions: group.Options,
			fiber.MethodTrace:   group.Trace,
		}
		for _, route := range router.Routes() {
			mapRoute, exists := m[route.Method]
			if exists {
				mapRoute(route.Path, route.Handler)
			}
		}
	}
	return nil
}
