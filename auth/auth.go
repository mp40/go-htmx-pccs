package auth

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/store"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	CountUserByEmail(email string) (int, error)
	GetUserByEmail(email string) (*store.User, error)
	AddUser(email string, hash string) (userID uuid.UUID, err error)
}

type Auth struct {
	store Store
	cost  int
}

func NewAuthService(store Store, cost int) *Auth {
	return &Auth{
		store: store,
		cost:  cost,
	}
}

func (a *Auth) SignIn(email string, password string) (*store.User, error) {
	user, err := a.store.GetUserByEmail(email)
	if err != nil {
		slog.Error("auth error getting user", "err", err)
		return nil, fmt.Errorf("500")
	}
	if user == nil {
		slog.Info("auth user not found", "email", email)
		// promote to typed error if/when needed
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		slog.Info("auth compare error")
		// promote to typed error if/when needed
		return nil, nil
	}

	return user, err
}

func (a *Auth) SignUp(email string, password string) (ID *uuid.UUID, err error) {
	userCount, err := a.store.CountUserByEmail(email)
	if err != nil {
		slog.Error("auth error counting by email", "err", err)
		return nil, fmt.Errorf("unexpected store error: %w", err)
	}
	if userCount != 0 {
		slog.Info("auth user email in use")
		return nil, fmt.Errorf("409")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), a.cost)
	if err != nil {
		slog.Error("auth error generating hash", "err", err)
		return nil, fmt.Errorf("unexpected store error: %w", err)
	}

	userID, err := a.store.AddUser(email, string(hash))
	if err != nil {
		slog.Error("auth error adding user", "err", err)
		return nil, fmt.Errorf("unexpected store error: %w", err)
	}

	return &userID, nil
}
