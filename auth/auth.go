package auth

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/data"
	"golang.org/x/crypto/bcrypt"
)

type Data interface {
	CountUserByEmail(email string) (int, error)
	GetUserByEmail(email string) (*data.User, error)
	AddUser(email string, hash string) (userID uuid.UUID, err error)
}

type Auth struct {
	data Data
}

func NewAuthService(data Data) *Auth {
	return &Auth{
		data: data,
	}
}

func (a *Auth) SignIn(email string, password string) (*data.User, error) {
	user, err := a.data.GetUserByEmail(email)
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
	userCount, err := a.data.CountUserByEmail(email)
	if err != nil {
		slog.Error("auth error counting by email", "err", err)
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}
	if userCount != 0 {
		slog.Info("auth user email in use")
		return nil, fmt.Errorf("409")
	}

	COST := os.Getenv("COST")
	if len(COST) == 0 {
		slog.Error("auth env COST not found")
		return nil, fmt.Errorf("500")
	}
	parsedCost, err := strconv.Atoi(COST)
	if err != nil {
		slog.Error("auth env COST parsing error", "err", err)
		return nil, fmt.Errorf("500")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), parsedCost)
	if err != nil {
		slog.Error("auth error generating hash", "err", err)
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}

	userID, err := a.data.AddUser(email, string(hash))
	if err != nil {
		slog.Error("auth error adding user", "err", err)
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}

	return &userID, nil
}
