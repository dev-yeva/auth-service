package storage

import (
	"auth/internal/model"
	"context"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (model.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int64) (model.App, error)
}
