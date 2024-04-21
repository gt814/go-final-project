package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go-final-project/config"
	"go-final-project/store"
	"net/http"
	"testing"
	"time"
)

func notFoundTask(t *testing.T, id string) {
	body, err := requestJSON("api/task?id="+id, nil, http.MethodGet)
	assert.NoError(t, err)
	var m map[string]any
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)
	_, ok := m["error"]
	assert.True(t, ok)
}

func TestDone(t *testing.T) {
	dbPath := config.GetDBFileTestPath()
	db, err := store.OpenDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	id := addTask(t, store.Task{
		Date:  now.Format(`20060102`),
		Title: "Свести баланс",
	})

	ret, err := postJSON("api/task/done?id="+id, nil, http.MethodPost)
	assert.NoError(t, err)
	assert.Empty(t, ret)
	notFoundTask(t, id)

	id = addTask(t, store.Task{
		Title:  "Проверить работу /api/task/done",
		Repeat: "d 3",
	})

	for i := 0; i < 3; i++ {
		ret, err := postJSON("api/task/done?id="+id, nil, http.MethodPost)
		assert.NoError(t, err)
		assert.Empty(t, ret)

		var task store.Task
		err = db.Get(&task, `SELECT * FROM scheduler WHERE id=?`, id)
		assert.NoError(t, err)
		now = now.AddDate(0, 0, 3)
		assert.Equal(t, task.Date, now.Format(`20060102`))
	}
}

func TestDelTask(t *testing.T) {
	dbPath := config.GetDBFileTestPath()
	db, err := store.OpenDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	id := addTask(t, store.Task{
		Title:  "Временная задача",
		Repeat: "d 3",
	})
	ret, err := postJSON("api/task?id="+id, nil, http.MethodDelete)
	assert.NoError(t, err)
	assert.Empty(t, ret)

	notFoundTask(t, id)

	ret, err = postJSON("api/task", nil, http.MethodDelete)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	ret, err = postJSON("api/task?id=wjhgese", nil, http.MethodDelete)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
}
