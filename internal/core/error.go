package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type ErrorHandler = func(*fiber.Ctx, error) error

func NewErrorHandler() ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var fibererr *fiber.Error
		if errors.As(err, &fibererr) {
			return ctx.Status(fibererr.Code).JSON(fiber.Map{"message": fibererr.Message})
		}
		return ctx.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
}
