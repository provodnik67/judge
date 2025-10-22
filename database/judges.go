package database

import (
	"database/sql"
	"time"
)

type Judge struct {
	ID          int
	Name        string
	Worldview   string
	Personality string
	Backstory   string
	IsActive    bool
	CreatedAt   time.Time
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

func GetAllJudges(db *sql.DB) ([]Judge, error) {
	rows, err := db.Query(`
		SELECT * FROM judges WHERE is_active = TRUE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var judges []Judge
	for rows.Next() {
		var j Judge
		err := rows.Scan(&j.ID, &j.Name, &j.Worldview, &j.Personality, &j.Backstory, &j.IsActive, &j.CreatedAt)
		if err != nil {
			return nil, err
		}
		judges = append(judges, j)
	}
	return judges, nil
}
