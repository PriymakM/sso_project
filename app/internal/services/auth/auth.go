package auth

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app ID")
	ErrUserExists         = errors.New("user does not exist")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (app models.App, err error)
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	
	return &Auth{usrSaver: userSaver, usrProvider: userProvider, log: log, appProvider: appProvider, tokenTTL: tokenTTL}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("username", email))
	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewTokens(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token", err)
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate passHash", err)
		return 0, err
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", err)

			return 0, fmt.Errorf("%s:%w", op, ErrUserExists)
		}
		log.Error("failed to save user", err)
		return 0, err
	}

	log.Info("user registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("user_id", userId))
	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("usert not found", err)

			return false, fmt.Errorf("%s:%w", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is admin", isAdmin))
	return isAdmin, nil
}
