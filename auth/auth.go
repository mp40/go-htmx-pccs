package auth

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Data interface {
	CountUserByEmail(email string) (int, error)
	AddUser(email string, hash string) (userID uuid.UUID, err error)
}

// temp in memory till need to impliment
type Auth struct {
	signedIn bool
	data     Data
}

func NewAuthService(data Data) *Auth {
	return &Auth{
		data: data,
	}
}

func (a *Auth) SignIn(email string, password string) error {
	a.signedIn = true
	return nil
}

func (a *Auth) SignUp(email string, password string) (ID *uuid.UUID, err error) {
	// validation

	userCount, err := a.data.CountUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("unexpected data error: %w", err)
	}
	if userCount != 0 {
		return nil, fmt.Errorf("409")
	}

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

	a.signedIn = true
	return &userID, nil
}
