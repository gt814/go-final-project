package store

import (
	"log"

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
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (s TaskStore) Create(t Task) (int64, error) {

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

func (s TaskStore) GetById(id int64) (Task, error) {
	task := Task{}
	err := s.db.Get(&task, `SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`, id)
	return task, err
}

func (s TaskStore) Update(t Task) error {
	_, err := s.db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`, t.Date, t.Title, t.Comment, t.Repeat, t.ID)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (s TaskStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}
