package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestUpsertEquipmentIDByCharacterID(t *testing.T) {
	t.Run("it returns error when load is not primary or secondary", func(t *testing.T) {
		store := NewStoreService(nil)

		err := store.UpsertEquipmentByCharacterID(1, uuid.New(), "xxxx", 1)

		if err == nil {
			t.Errorf("expect error, got nil")
		}
	})
}
