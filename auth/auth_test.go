package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/data"
	"golang.org/x/crypto/bcrypt"
)

type StubData struct {
	newID         uuid.UUID
	user          *data.User
	count         int
	countErr      error
	err           error
	spyGetByEmail int
}

func (d *StubData) CountUserByEmail(email string) (int, error) {
	return d.count, d.countErr
}

func (d *StubData) GetUserByEmail(email string) (user *data.User, err error) {
	d.spyGetByEmail++
	return d.user, d.err
}

func (d *StubData) AddUser(email string, hash string) (userID uuid.UUID, err error) {
	return d.newID, d.err
}

func TestSignUp(t *testing.T) {
	t.Setenv("SALT", "1")
	t.Run("it should return pointer to ID when email not in use", func(t *testing.T) {
		newID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")
		stubData := StubData{}
		stubData.newID = newID
		auth := NewAuthService(&stubData)

		got, err := auth.SignUp("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if *got != stubData.newID {
			t.Errorf("got %v, want %v", *got, stubData.newID)
		}
	})

	t.Run("it should return conflict error when email in use", func(t *testing.T) {
		stubData := StubData{}
		stubData.count = 1
		auth := NewAuthService(&stubData)

		_, err := auth.SignUp("fake@email.com", "fake-password")

		if err == nil {
			t.Errorf("expected error got nil")
		}

		if err.Error() != "409" {
			t.Errorf("got expected 409 error got %v", err.Error())
		}
	})

	t.Run("it should return error on data error", func(t *testing.T) {
		stubData := StubData{}
		stubData.countErr = fmt.Errorf("fake error")
		auth := NewAuthService(&stubData)

		_, err := auth.SignUp("fake@email.com", "fake-password")

		if err == nil {
			t.Errorf("expected error got nil")
		}

		if !strings.Contains(err.Error(), "unexpected data error") {
			t.Errorf("got unexpected error got %v", err.Error())
		}
	})
}

func TestSignIn(t *testing.T) {
	password := "fake-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		panic("test fixture failure")
	}
	t.Run("it should return pointer to User when user found by by email and hash", func(t *testing.T) {
		stubData := StubData{}
		user := data.User{Hash: string(hash)}
		stubData.user = &user
		auth := NewAuthService(&stubData)

		got, err := auth.SignIn("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if stubData.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to data, got %v, want 1", stubData.spyGetByEmail)
		}

		if got == nil {
			t.Errorf("got nil, want user")
		}
	})

	t.Run("it should return nil when user not found by email", func(t *testing.T) {
		stubData := StubData{}
		auth := NewAuthService(&stubData)

		got, err := auth.SignIn("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if stubData.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to data, got %v, want 1", stubData.spyGetByEmail)
		}

		if got != nil {
			t.Errorf("got user, want nil")
		}
	})

	t.Run("it should return pointer to nil when provided password does not match", func(t *testing.T) {
		stubData := StubData{}
		user := data.User{Hash: "something-else"}
		stubData.user = &user
		auth := NewAuthService(&stubData)

		got, err := auth.SignIn("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if stubData.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to data, got %v, want 1", stubData.spyGetByEmail)
		}

		if got != nil {
			t.Errorf("expected nil, got user")
		}
	})

	t.Run("it should return error on data error", func(t *testing.T) {
		stubData := StubData{err: fmt.Errorf("fake")}
		auth := NewAuthService(&stubData)

		_, err := auth.SignIn("fake@email.com", "fake-password")

		if stubData.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to data, got %v, want 1", stubData.spyGetByEmail)
		}

		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
