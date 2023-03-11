package store

import (
	"BrainTApp/internal/model/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	Find(int) (*entity.User, error)
	FindByEmail(string) (*entity.User, error)
}

type ResultRepository interface {
	CreateShulte(result *entity.Result) error
	Find(user_id int, test_id int) (*entity.Result, error)
	FindAll(user_id int) ([]entity.Result, error)
	CreateArithmetic(result *entity.Result) error
	CreateMemorize(result *entity.Result) error
	TestSetWords(num int, speech int) ([]string, error)
}
