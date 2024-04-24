package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-final-project/api"
	"go-final-project/config"
	"go-final-project/store"
	"net/http"
	"testing"
	"time"
)

func addTask(t *testing.T, task store.Task) string {
	ret, err := postJSON("api/task", map[string]any{
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}, http.MethodPost)
	assert.NoError(t, err)
	assert.NotNil(t, ret["id"])
	id := fmt.Sprint(ret["id"])
	assert.NotEmpty(t, id)
	return id
}

func getTasks(t *testing.T, search string) []store.Task {
	url := "api/tasks"
	if config.GetSearch() {
		url += "?search=" + search
	}
	body, err := requestJSON(url, nil, http.MethodGet)
	assert.NoError(t, err)

	var tasksResponse api.TasksResponse
	err = json.Unmarshal(body, &tasksResponse)
	assert.NoError(t, err)
	return tasksResponse.Tasks
}

func TestTasks(t *testing.T) {
	dbPath := config.GetDBFileTestPath()
	db, err := store.OpenDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	_, err = db.Exec("DELETE FROM scheduler")
	assert.NoError(t, err)

	tasks := getTasks(t, "")
	assert.NotNil(t, tasks)
	assert.Empty(t, tasks)

	addTask(t, store.Task{
		Date:    now.Format(`20060102`),
		Title:   "Просмотр фильма",
		Comment: "с попкорном",
		Repeat:  "",
	})
	now = now.AddDate(0, 0, 1)
	date := now.Format(`20060102`)
	addTask(t, store.Task{
		Date:    date,
		Title:   "Сходить в бассейн",
		Comment: "",
		Repeat:  "",
	})
	addTask(t, store.Task{
		Date:    date,
		Title:   "Оплатить коммуналку",
		Comment: "",
		Repeat:  "d 30",
	})
	tasks = getTasks(t, "")
	assert.Equal(t, len(tasks), 3)

	now = now.AddDate(0, 0, 2)
	date = now.Format(`20060102`)
	addTask(t, store.Task{
		Date:    date,
		Title:   "Поплавать",
		Comment: "Бассейн с тренером",
		Repeat:  "d 7",
	})
	addTask(t, store.Task{
		Date:    date,
		Title:   "Позвонить в УК",
		Comment: "Разобраться с горячей водой",
		Repeat:  "",
	})
	addTask(t, store.Task{
		Date:    date,
		Title:   "Встретится с Васей",
		Comment: "в 18:00",
		Repeat:  "",
	})

	tasks = getTasks(t, "")
	assert.Equal(t, len(tasks), 6)

	if !config.GetSearch() {
		return
	}
	tasks = getTasks(t, "УК")
	assert.Equal(t, len(tasks), 1)
	tasks = getTasks(t, now.Format(`02.01.2006`))
	assert.Equal(t, len(tasks), 3)

}
