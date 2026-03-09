package auth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type StubData struct {
	newID    uuid.UUID
	count    int
	countErr error
	addErr   error
}

func (d *StubData) CountUserByEmail(email string) (int, error) {
	return d.count, d.countErr
}

func (d *StubData) AddUser(email string, hash string) (userID uuid.UUID, err error) {
	return d.newID, d.addErr
}

func TestSignUp(t *testing.T) {
	t.Run("it should return nil error when email not in use", func(t *testing.T) {
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
