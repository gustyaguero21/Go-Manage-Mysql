package services

import (
	"context"
	"go-manage-mysql/internal/models"
)

type UserServices interface {
	CreateUser(ctx context.Context, user models.User) (created models.User, err error)
	SearchUser(ctx context.Context, username string) (user models.User, err error)
}
