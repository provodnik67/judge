package database

import (
	"database/sql"
)

type Judge struct {
	ID          int
	Name        string
	Worldview   string
	Personality string
	Backstory   string
	IsActive    bool
}

func CreateJudge(db *sql.DB, judge Judge) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO judges (name, worldview, personality_prompt, backstory, is_active) 
		VALUES (?, ?, ?, ?, ?)`,
		judge.Name, judge.Worldview, judge.Personality, judge.Backstory, judge.IsActive,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
