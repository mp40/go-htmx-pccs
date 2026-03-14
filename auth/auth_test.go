package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/store"
	"golang.org/x/crypto/bcrypt"
)

type StubStore struct {
	newID         uuid.UUID
	user          *store.User
	count         int
	countErr      error
	err           error
	spyGetByEmail int
}

func (d *StubStore) CountUserByEmail(email string) (int, error) {
	return d.count, d.countErr
}

func (d *StubStore) GetUserByEmail(email string) (user *store.User, err error) {
	d.spyGetByEmail++
	return d.user, d.err
}

func (d *StubStore) AddUser(email string, hash string) (userID uuid.UUID, err error) {
	return d.newID, d.err
}

func TestSignUp(t *testing.T) {
	t.Setenv("SALT", "1")
	t.Run("it should return pointer to ID when email not in use", func(t *testing.T) {
		newID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")
		stubStore := StubStore{}
		stubStore.newID = newID
		auth := NewAuthService(&stubStore)

		got, err := auth.SignUp("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if *got != stubStore.newID {
			t.Errorf("got %v, want %v", *got, stubStore.newID)
		}
	})

	t.Run("it should return conflict error when email in use", func(t *testing.T) {
		stubStore := StubStore{}
		stubStore.count = 1
		auth := NewAuthService(&stubStore)

		_, err := auth.SignUp("fake@email.com", "fake-password")

		if err == nil {
			t.Errorf("expected error got nil")
		}

		if err.Error() != "409" {
			t.Errorf("got expected 409 error got %v", err.Error())
		}
	})

	t.Run("it should return error on store error", func(t *testing.T) {
		stubStore := StubStore{}
		stubStore.countErr = fmt.Errorf("fake error")
		auth := NewAuthService(&stubStore)

		_, err := auth.SignUp("fake@email.com", "fake-password")

		if err == nil {
			t.Errorf("expected error got nil")
		}

		if !strings.Contains(err.Error(), "unexpected store error") {
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
		stubStore := StubStore{}
		user := store.User{Hash: string(hash)}
		stubStore.user = &user
		auth := NewAuthService(&stubStore)

		got, err := auth.SignIn("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if stubStore.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to store, got %v, want 1", stubStore.spyGetByEmail)
		}

		if got == nil {
			t.Errorf("got nil, want user")
		}
	})

	t.Run("it should return nil when user not found by email", func(t *testing.T) {
		stubStore := StubStore{}
		auth := NewAuthService(&stubStore)

		got, err := auth.SignIn("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if stubStore.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to store, got %v, want 1", stubStore.spyGetByEmail)
		}

		if got != nil {
			t.Errorf("got user, want nil")
		}
	})

	t.Run("it should return pointer to nil when provided password does not match", func(t *testing.T) {
		stubStore := StubStore{}
		user := store.User{Hash: "something-else"}
		stubStore.user = &user
		auth := NewAuthService(&stubStore)

		got, err := auth.SignIn("fake@email.com", "fake-password")
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
		}

		if stubStore.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to store, got %v, want 1", stubStore.spyGetByEmail)
		}

		if got != nil {
			t.Errorf("expected nil, got user")
		}
	})

	t.Run("it should return error on store error", func(t *testing.T) {
		stubStore := StubStore{err: fmt.Errorf("fake")}
		auth := NewAuthService(&stubStore)

		_, err := auth.SignIn("fake@email.com", "fake-password")

		if stubStore.spyGetByEmail != 1 {
			t.Errorf("unexpected calls to store, got %v, want 1", stubStore.spyGetByEmail)
		}

		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
