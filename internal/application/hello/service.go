package hello

type Service interface {
	Hello() string
	World() string
}

type helloService struct {
	repository Repository
}

func NewHelloService(repository Repository) Service {
	return helloService{
		repository: repository,
	}
}

func (h helloService) Hello() string {
	return "World"
}

func (h helloService) World() string {
	return "Hello"
}
