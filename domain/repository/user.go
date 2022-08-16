package repository

import (
	"context"

	"github.com/ryoh07/gin-clean-webapp/domain/entity"
)

type IUserRepository interface {
	GetUser(userId string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	CreateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userId string) error
}
