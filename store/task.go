package store

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"time"
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

func (s TaskStore) GetTaskList(count int) ([]Task, error) {
	var res []Task

	query := fmt.Sprintf("SELECT id, date, title, comment, repeat FROM scheduler LIMIT %d", count)
	rows, err := s.db.Query(query)
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
	return err
}

func (s TaskStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	return err
}

func NextDate(now time.Time, strDate string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule not specified")
	}

	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return "", errors.New("invalid date format")
	}

	repeatFields := strings.Fields(repeat)
	if len(repeatFields) < 1 {
		return "", errors.New("invalid repeat format")
	}

	rule := repeatFields[0]
	repeatValue := ""
	if len(repeatFields) > 1 {
		repeatValue = repeatFields[1]
	}

	nextDate := date
	switch rule {
	case "d":
		if repeatValue == "" {
			return "", errors.New("invalid repeat format")
		}
		days, err := strconv.Atoi(repeatValue)
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid repeat format")
		}
		nextDate = nextDate.AddDate(0, 0, days)
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	case "w":
		if repeatValue == "" {
			return "", errors.New("invalid repeat format")
		}
		targetWeekdays := make(map[time.Weekday]bool)
		for _, dayStr := range strings.Split(repeatValue, ",") {
			dayInt, err := strconv.Atoi(dayStr)
			if err != nil || dayInt < 1 || dayInt > 7 {
				return "", errors.New("invalid weekday format")
			}
			if dayInt == 7 {
				dayInt = 0
			}
			targetWeekdays[time.Weekday(dayInt)] = true
		}

		nextDate = nextDate.AddDate(0, 0, 1)
		for {
			if now.Before(nextDate) && targetWeekdays[nextDate.Weekday()] {
				break
			}
			nextDate = nextDate.AddDate(0, 0, 1)
		}

	default:
		return "", errors.New("invalid repeat rule")
	}

	return nextDate.Format("20060102"), nil
}
