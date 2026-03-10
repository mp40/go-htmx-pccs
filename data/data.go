package data

import (
	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID `json:"id"`
	Hash string    `json:"hash"`
}

type Data struct{}

func NewDataService() *Data {
	return &Data{}
}

func (d *Data) CountUserByEmail(email string) (int, error) {
	return 0, nil
}

func (d *Data) GetUserByEmail(email string) (*User, error) {
	user := User{}
	return &user, nil
}

func (d *Data) AddUser(email string, hash string) (id uuid.UUID, err error) {
	ID := uuid.New()
	return ID, nil
}
