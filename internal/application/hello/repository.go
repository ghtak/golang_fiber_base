package hello

import (
	"github.com/golang_fiber_base/internal/core"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
}

type helloRepository struct {
	database *sqlx.DB
}

func NewHelloRepository(p core.Param) Repository {
	return helloRepository{
		database: p.Database,
	}
}
