package service

import (
	"errors"
	"fmt"
	"go-final-project/config"
	"go-final-project/store"
	"log"
	"strconv"
	"strings"
	"time"
)

var taskStore store.TaskStore

func SetTaskStore(ts store.TaskStore) {
	taskStore = ts
}

func CheckTask(t store.Task) (store.Task, error) {
	if t.Title == "" {
		return t, errors.New("task title is not specified")
	}

	if t.Date != "" {
		_, err := time.Parse("20060102", t.Date)

		if err != nil {
			return t, errors.New("invalid date format")
		}
	}
	return t, nil
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

func Create(task store.Task) (string, error) {
	task, err := enrichTask(task)
	if err != nil {
		return "", err
	}

	id, err := taskStore.Create(task)

	if err != nil {
		log.Println("Create err=", err.Error())
		return "", err
	}

	return fmt.Sprint(id), nil
}

func GetTasks() ([]store.Task, error) {
	var tasks []store.Task

	count := config.GetTaskLimit()
	tasks, err := taskStore.GetTaskList(count)

	if err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []store.Task{}
	}

	return tasks, nil
}

func GetById(id int64) (store.Task, error) {
	return taskStore.GetById(id)
}

func Update(task store.Task) error {
	task, err := enrichTask(task)
	if err != nil {
		return err
	}
	return taskStore.Update(task)
}

func Done(t store.Task) error {
	id, err := strconv.ParseInt(t.ID, 10, 64)
	if err != nil {
		return err
	}

	if t.Repeat == "" {
		err = Delete(id)

		if err != nil {
			log.Println("Delete err=", err.Error())
			return err
		}
	} else {
		t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return err
		}

		err = Update(t)
		if err != nil {
			log.Println("Update err=", err.Error())
			return err
		}
	}
	return nil
}

func Delete(id int64) error {
	return taskStore.Delete(id)
}

func enrichTask(t store.Task) (store.Task, error) {
	var taskDate time.Time
	var err error
	if t.Date == "" {
		taskDate = time.Now()
		t.Date = taskDate.Format("20060102")
	} else {
		taskDate, err = time.Parse("20060102", t.Date)

		if err != nil {
			return t, errors.New("invalid date format")
		}
	}

	if t.Repeat != "" {
		if taskDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
			t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)

			if err != nil {
				return t, errors.New("invalid date format")
			}
		}
	} else {
		if taskDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
			t.Date = time.Now().Format("20060102")
		}
	}

	return t, nil
}
