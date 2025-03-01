package services

import (
	"context"
	"go-manage-mysql/internal/models"
)

type UserServices interface {
	CreateUser(ctx context.Context, user models.User) (created models.User, err error)
	SearchUser(ctx context.Context, username string) (user models.User, err error)
	UpdateUser(ctx context.Context, username string, update models.User) (err error)
	DeleteUser(ctx context.Context, username string) (err error)
	ChangeUserPwd(ctx context.Context, username string, newPwd string) (err error)
}
