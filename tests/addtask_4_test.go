package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-final-project/config"
	"go-final-project/store"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func requestJSON(apipath string, values map[string]any, method string) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	if len(values) > 0 {
		data, err = json.Marshal(values)
		if err != nil {
			return nil, err
		}
	}
	var resp *http.Response

	req, err := http.NewRequest(method, getURL(apipath), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if len(config.GetToken()) > 0 {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		jar.SetCookies(req.URL, []*http.Cookie{
			{
				Name:  "token",
				Value: config.GetToken(),
			},
		})
		client.Jar = jar
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return io.ReadAll(resp.Body)
}

func postJSON(apipath string, values map[string]any, method string) (map[string]any, error) {
	var (
		m   map[string]any
		err error
	)

	body, err := requestJSON(apipath, values, method)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &m)
	return m, err
}

func TestAddTask(t *testing.T) {
	dbPath := config.GetDBFileTestPath()
	db, err := store.OpenDB(dbPath)
	assert.NoError(t, err)
	defer db.Close()

	tbl := []store.Task{
		{Date: "20240129", Title: "", Comment: "", Repeat: ""},
		{Date: "20240192", Title: "Qwerty", Comment: "", Repeat: ""},
		{Date: "28.01.2024", Title: "Заголовок", Comment: "", Repeat: ""},
		{Date: "20240112", Title: "Заголовок", Comment: "", Repeat: "w"},
		{Date: "20240212", Title: "Заголовок", Comment: "", Repeat: "ooops"},
	}
	for _, v := range tbl {
		m, err := postJSON("api/task", map[string]any{
			"date":    v.Date,
			"title":   v.Title,
			"comment": v.Comment,
			"repeat":  v.Repeat,
		}, http.MethodPost)
		assert.NoError(t, err)

		e, ok := m["error"]
		assert.False(t, !ok || len(fmt.Sprint(e)) == 0,
			"Ожидается ошибка для задачи %v", v)
	}

	now := time.Now()

	check := func() {
		for _, v := range tbl {
			today := v.Date == "today"
			if today {
				v.Date = now.Format(`20060102`)
			}
			log.Println("v=", v) //delete this print before merge
			m, err := postJSON("api/task", map[string]any{
				"date":    v.Date,
				"title":   v.Title,
				"comment": v.Comment,
				"repeat":  v.Repeat,
			}, http.MethodPost)
			assert.NoError(t, err)

			e, ok := m["error"]
			if ok && len(fmt.Sprint(e)) > 0 {
				t.Errorf("Неожиданная ошибка %v для задачи %v", e, v)
				continue
			}
			var task store.Task
			var mid any
			mid, ok = m["id"]
			if !ok {
				t.Errorf("Не возвращён id для задачи %v", v)
				continue
			}
			id := fmt.Sprint(mid)

			err = db.Get(&task, `SELECT * FROM scheduler WHERE id=?`, id)
			log.Println("id=", id)     //delete this print before merge
			log.Println("task=", task) //delete this print before merge
			assert.NoError(t, err)
			assert.Equal(t, id, strconv.FormatInt(task.ID, 10))

			assert.Equal(t, v.Title, task.Title)
			assert.Equal(t, v.Comment, task.Comment)
			assert.Equal(t, v.Repeat, task.Repeat)
			if task.Date < now.Format(`20060102`) {
				t.Errorf("Дата не может быть меньше сегодняшней %v", v)
				continue
			}
			if today && task.Date != now.Format(`20060102`) {
				t.Errorf("Дата должна быть сегодняшняя %v", v)
			}
		}
	}

	tbl = []store.Task{
		{Date: "", Title: "Заголовок", Comment: "", Repeat: ""},
		{Date: "20231220", Title: "Сделать что-нибудь", Comment: "Хорошо отдохнуть", Repeat: ""},
		{Date: "20240108", Title: "Уроки", Comment: "", Repeat: "d 10"},
		{Date: "20240102", Title: "Отдых в Сочи", Comment: "На лыжах", Repeat: "y"},
		{Date: "today", Title: "Фитнес", Comment: "", Repeat: "d 1"},
		{Date: "today", Title: "Шмитнес", Comment: "", Repeat: ""},
	}
	check()
	if config.GetFullNextDate() {
		tbl = []store.Task{
			{Date: "20240129", Title: "Сходить в магазин", Comment: "", Repeat: "w 1,3,5"},
		}
		check()
	}
}
