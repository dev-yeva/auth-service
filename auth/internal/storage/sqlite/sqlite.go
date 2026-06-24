package sqlite

import (
	"auth/internal/model"
	"auth/internal/storage"
	"auth/lib"
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) *Storage {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO users(email, password_hash, is_admin) 
		VALUES (?, ?, ?)`,
		email, passwordHash, false,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, storage.ErrUserExists
		}
		return 0, lib.ErrWrap("failed to save user data", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, lib.ErrWrap("failed to get saved user id", err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (model.User, error) {
	user := model.User{}

	row := s.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, is_admin 
		FROM users 
		WHERE email = ?`,
		email,
	)

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, storage.ErrUserNotFound
		}
		return model.User{}, lib.ErrWrap("failed to extract user data", err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	var isAdmin bool

	row := s.db.QueryRowContext(ctx, `
		SELECT is_admin 
		FROM users 
		WHERE id = ?`,
		userId,
	)

	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrUserNotFound
		}
		return false, lib.ErrWrap("failed to check admin status", err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appId int64) (model.App, error) {
	app := model.App{}

	row := s.db.QueryRowContext(ctx, `
		SELECT id, title, secret_key
		FROM apps 
		WHERE id = ?`,
		appId,
	)

	err := row.Scan(&app.ID, &app.Title, &app.SecretKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.App{}, storage.ErrAppNotFound
		}
		return model.App{}, lib.ErrWrap("failed to extract app data", err)
	}

	return app, nil
}
