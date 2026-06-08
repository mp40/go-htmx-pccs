package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Hash      string    `json:"hash"`
	CreatedAt int64     `json:"created_at"` // unix timestamp
	UpdatedAt int64     `json:"updated_at"` // unix timestamp
}

type Character struct {
	ID                       uuid.UUID `json:"id"`
	UserID                   uuid.UUID `json:"user_id"`
	Name                     string    `json:"name"`
	Str                      int       `json:"str"`
	Int                      int       `json:"int"`
	Wil                      int       `json:"wil"`
	Hlt                      int       `json:"hlt"`
	Agi                      int       `json:"agi"`
	Tch                      int       `json:"tch"`
	GunCombatLearningPoints  float32   `json:"gun_combat_learning_points"`
	HandToHandLearningPoints float32   `json:"hand_to_hand_learning_points"`
	CreatedAt                int64     `json:"created_at"` // unix timestamp
	UpdatedAt                int64     `json:"updated_at"` // unix timestamp
}

type Store struct {
	db *sql.DB
}

func NewStoreService(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CountUserByEmail(email string) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return count, err
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	user := User{}
	var idStr string
	err := s.db.QueryRow("SELECT id, email, hash, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&idStr, &user.Email, &user.Hash, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CountUserByID(ID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", ID.String()).Scan(&count)
	return count, err
}

func (s *Store) AddUser(email string, hash string) (uuid.UUID, error) {
	ID := uuid.New()
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO users (id, email, hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		ID.String(), email, hash, now, now,
	)
	if err != nil {
		return uuid.Nil, err
	}
	return ID, nil
}

func (s *Store) AddCharacter(character Character) (*Character, error) {
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO characters (id, user_id, name, str, int, wil, hlt, agi, tch, gun_combat_learning_points, hand_to_hand_learning_points, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		character.ID, character.UserID, character.Name, character.Str, character.Int, character.Wil, character.Hlt, character.Agi, character.Tch, character.GunCombatLearningPoints, character.HandToHandLearningPoints, now, now,
	)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (s *Store) GetCharactersByUserID(userID uuid.UUID) ([]Character, error) {
	rows, err := s.db.Query("SELECT * FROM characters")
	if err != nil {
		return []Character{}, err
	}

	var characters []Character

	for rows.Next() {
		character := Character{}
		var idStr string
		var userIdStr string
		err := rows.Scan(&idStr, &userIdStr, &character.Name, &character.Str, &character.Int, &character.Wil, &character.Hlt, &character.Agi, &character.CreatedAt, &character.UpdatedAt)
		if err != nil {
			return []Character{}, err
		}

		character.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		character.UserID, err = uuid.Parse(userIdStr)
		if err != nil {
			return nil, err
		}

		characters = append(characters, character)
	}

	return characters, nil
}

func (s *Store) GetCharacterByID(ID uuid.UUID) (*Character, error) {
	character := Character{}
	var idStr string
	var userIdStr string
	err := s.db.QueryRow("SELECT id, user_id, name, str, int, wil, hlt, agi, created_at, updated_at FROM characters WHERE id = ?", ID).
		Scan(&idStr, &userIdStr, &character.Name, &character.Str, &character.Int, &character.Wil, &character.Hlt, &character.Agi, &character.CreatedAt, &character.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	character.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	character.UserID, err = uuid.Parse(userIdStr)
	if err != nil {
		return nil, err
	}
	return &character, nil
}
