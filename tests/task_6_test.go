package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-final-project/config"
	"go-final-project/store"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	dbPath := config.GetDBFileTestPath()
	db, err := store.OpenDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()

	task := store.Task{
		Date:    now.Format(`20060102`),
		Title:   "Созвон в 16:00",
		Comment: "Обсуждение планов",
		Repeat:  "d 5",
	}

	todo := addTask(t, task)

	body, err := requestJSON("api/task", nil, http.MethodGet)
	assert.NoError(t, err)
	var responseTask store.Task
	err = json.Unmarshal(body, &responseTask)

	assert.False(t, len(fmt.Sprint(err)) == 0, "Ожидается ошибка для вызова /api/task")

	body, err = requestJSON("api/task?id="+todo, nil, http.MethodGet)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &responseTask)
	assert.NoError(t, err)

	taskId, _ := strconv.ParseInt(todo, 10, 64)
	assert.Equal(t, taskId, responseTask.ID)
	assert.Equal(t, task.Date, responseTask.Date)
	assert.Equal(t, task.Title, responseTask.Title)
	assert.Equal(t, task.Comment, responseTask.Comment)
	assert.Equal(t, task.Repeat, responseTask.Repeat)
}

type fulltask struct {
	id string
	store.Task
}

func TestEditTask(t *testing.T) {
	dbPath := config.GetDBFileTestPath()
	db, err := store.OpenDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()

	tsk := store.Task{
		Date:    now.Format(`20060102`),
		Title:   "Заказать пиццу",
		Comment: "в 17:00",
		Repeat:  "",
	}

	id := addTask(t, tsk)

	tbl := []fulltask{
		{"", store.Task{Date: "20240129", Title: "Тест", Comment: "", Repeat: ""}},
		{"abc", store.Task{Date: "20240129", Title: "Тест", Comment: "", Repeat: ""}},
		{"7645346343", store.Task{Date: "20240129", Title: "Тест", Comment: "", Repeat: ""}},
	}
	for _, v := range tbl {
		m, err := postJSON("api/task", map[string]any{
			"id":      v.id,
			"date":    v.Date,
			"title":   v.Title,
			"comment": v.Comment,
			"repeat":  v.Repeat,
		}, http.MethodPut)
		assert.NoError(t, err)

		var errVal string
		e, ok := m["error"]
		if ok {
			errVal = fmt.Sprint(e)
		}
		assert.NotEqual(t, len(errVal), 0, "Ожидается ошибка для значения %v", v)
	}

	idParam, _ := strconv.ParseInt(id, 10, 64)
	tbl2 := []store.Task{
		{ID: idParam, Date: "20240129", Title: "", Comment: "", Repeat: ""},
		{ID: idParam, Date: "20240192", Title: "Qwerty", Comment: "", Repeat: ""},
		{ID: idParam, Date: "28.01.2024", Title: "Заголовок", Comment: "", Repeat: ""},
		{ID: idParam, Date: "20240212", Title: "Заголовок", Comment: "", Repeat: "ooops"},
	}
	for _, v := range tbl2 {
		m, err := postJSON("api/task", map[string]any{
			"id":      v.ID,
			"date":    v.Date,
			"title":   v.Title,
			"comment": v.Comment,
			"repeat":  v.Repeat,
		}, http.MethodPut)
		assert.NoError(t, err)

		var errVal string
		e, ok := m["error"]
		if ok {
			errVal = fmt.Sprint(e)
		}
		assert.NotEqual(t, len(errVal), 0, "Ожидается ошибка для значения %v", v)
	}

	updateTask := func(newVals map[string]any) {
		mupd, err := postJSON("api/task", newVals, http.MethodPut)
		assert.NoError(t, err)

		e, ok := mupd["error"]
		assert.False(t, ok && fmt.Sprint(e) != "")

		var task store.Task
		err = db.Get(&task, `SELECT * FROM scheduler WHERE id=?`, id)
		assert.NoError(t, err)

		assert.Equal(t, id, strconv.FormatInt(task.ID, 10))
		assert.Equal(t, newVals["title"], task.Title)
		if _, is := newVals["comment"]; !is {
			newVals["comment"] = ""
		}
		if _, is := newVals["repeat"]; !is {
			newVals["repeat"] = ""
		}
		assert.Equal(t, newVals["comment"], task.Comment)
		assert.Equal(t, newVals["repeat"], task.Repeat)
		now := time.Now().Format(`20060102`)
		if task.Date < now {
			t.Errorf("Дата не может быть меньше сегодняшней")
		}
	}

	updateTask(map[string]any{
		"id":      idParam,
		"date":    now.Format(`20060102`),
		"title":   "Заказать хинкали",
		"comment": "в 18:00",
		"repeat":  "d 7",
	})
}
