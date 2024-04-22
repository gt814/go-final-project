package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type TaskStore struct {
	db *sqlx.DB
}

func NewTaskStore(db *sqlx.DB) TaskStore {
	return TaskStore{db: db}
}

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (s TaskStore) Add(t Task) (int64, error) {

	res, err := s.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
	VALUES (?, ?, ?, ?)`, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

func (s TaskStore) GetAll() ([]Task, error) {
	var res []Task

	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Comment)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}
