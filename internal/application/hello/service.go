package hello

import "github.com/gofiber/fiber/v2"

type Service interface {
	Hello() string
	World() string
	Error() error
}

type helloService struct {
	repository Repository
}

func NewHelloService(repository Repository) Service {
	return &helloService{
		repository: repository,
	}
}

func (h *helloService) Hello() string {
	return "World"
}

func (h *helloService) World() string {
	return "Hello"
}

func (h *helloService) Error() error {
	return &fiber.Error{
		Code:    400,
		Message: "message",
	}
}
