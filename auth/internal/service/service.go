package service

import (
	"auth/internal/jwt"
	"auth/internal/storage"
	"auth/lib"
	"context"
	"errors"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	logger       *slog.Logger
	userSaver    storage.UserSaver
	userProvider storage.UserProvider
	appProvider  storage.AppProvider
	tokenTTL     time.Duration
}

func (a *AuthService) Register(ctx context.Context, email, password string) (int64, error) {

	logger := a.logger.With("email", email, "op", "service.Register")
	logger.Info("registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Warn("registration error", "err", err)
		return 0, lib.ErrWrap("failed to generate password hash", err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passwordHash)
	if err != nil {
		logger.Warn("registration error", "err", err)
		if errors.Is(err, storage.ErrUserExists) {
			return 0, ErrUserExists
		}
		return 0, err
	}

	logger.Info("successful registration")
	return id, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string, appId int64) (string, error) {
	logger := a.logger.With("email", email, "op", "service.Login")
	logger.Info("logging in user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		logger.Warn("login error", "err", err)
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", lib.ErrWrap("failed to get user", err)
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		logger.Warn("login error", "err", err)
		return "", ErrInvalidPassword
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return "", ErrAppNotFound
		}
		logger.Warn("login error", "err", err)
		return "", lib.ErrWrap("failed to get app", err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		logger.Warn("login error", "err", err)
		return "", lib.ErrWrap("failed to generate token", err)
	}

	logger.Info("successful login")

	return token, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	logger := a.logger.With("userId", userId, "op", "service.IsAdmin")
	logger.Info("checking admin status")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		logger.Warn("failed to check admin status", "err", err)
		if errors.Is(err, storage.ErrUserNotFound) {
			return false, ErrUserNotFound
		}
		return false, err
	}
	return isAdmin, nil
}

func New(
	logger *slog.Logger,
	userSaver storage.UserSaver,
	userProvider storage.UserProvider,
	appProvider storage.AppProvider,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		logger:       logger,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}
