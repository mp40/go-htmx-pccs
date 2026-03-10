package auth

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/data"
	"golang.org/x/crypto/bcrypt"
)

type Data interface {
	CountUserByEmail(email string) (int, error)
	GetUserByEmail(email string) (*data.User, error)
	GetUserByEmailAndHash(email string, hash string) (*data.User, error)
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
	// bad naming
	SALT := os.Getenv("SALT")
	if len(SALT) == 0 {
		return nil, fmt.Errorf("500")
	}
	parsedSalt, err := strconv.Atoi(SALT)
	if err != nil {
		return nil, fmt.Errorf("500")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), parsedSalt)
	if err != nil {
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}

	user, err := a.data.GetUserByEmailAndHash(email, string(hash))

	return user, err
}

func (a *Auth) SignUp(email string, password string) (ID *uuid.UUID, err error) {
	userCount, err := a.data.CountUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}
	if userCount != 0 {
		return nil, fmt.Errorf("409")
	}

	// bad naming
	SALT := os.Getenv("SALT")
	if len(SALT) == 0 {
		return nil, fmt.Errorf("500")
	}
	parsedSalt, err := strconv.Atoi(SALT)
	if err != nil {
		return nil, fmt.Errorf("500")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), parsedSalt)
	if err != nil {
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}

	userID, err := a.data.AddUser(email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}

	return &userID, nil
}
