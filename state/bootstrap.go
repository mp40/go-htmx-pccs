package state

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed seeds/*.sql
var seedFS embed.FS

func Bootstrap(db *sql.DB) error {
	entries, err := seedFS.ReadDir("seeds")
	if err != nil {
		return fmt.Errorf("read seeds dir: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	for _, entry := range entries {
		path := "seeds/" + entry.Name()
		content, err := seedFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if _, err := tx.Exec(string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", path, err)
		}
	}

	return tx.Commit()
}
