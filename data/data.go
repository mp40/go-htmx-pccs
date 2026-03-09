package data

import (
	"github.com/google/uuid"
)

type Data struct{}

func NewDataService() *Data {
	return &Data{}
}

func (d *Data) CountUserByEmail(email string) (int, error) {
	return 1, nil
}

func (d *Data) AddUser(email string, hash string) (id *uuid.UUID, err error) {
	ID := uuid.New()
	return &ID, nil
}
