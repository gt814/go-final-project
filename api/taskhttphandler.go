package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-final-project/service"
	"go-final-project/store"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TaskIdResponse struct {
	ID string `json:"id"`
}

type TasksResponse struct {
	Tasks []store.Task `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Response struct{}

var emptyResponse = Response{}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "Invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("20060102", dateParam)
	if err != nil {
		http.Error(w, "Invalid 'date' parameter", http.StatusBadRequest)
		return
	}

	nextDate, err := service.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task store.Task
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("read body, err=%w", err)
		fmt.Println(err)
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf("unmarshal task, err=%w", err)
		fmt.Println(err)
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	task, err = checkTask(task)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := service.Create(task)

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, TaskIdResponse{ID: id}, http.StatusCreated)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := service.GetTasks()

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, TasksResponse{Tasks: tasks}, http.StatusOK)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task ID not specified"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "Invalid task ID"}, http.StatusBadRequest)
		return
	}

	task, err := service.GetById(id)

	if task.ID == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
		return
	}

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, task, http.StatusOK)
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task store.Task
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = fmt.Errorf("read body, err=%w", err)
		fmt.Println(err)
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = fmt.Errorf("unmarshal task, err=%w", err)
		fmt.Println(err)
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		makeHttpResponse(w, ErrorResponse{Error: "ID not specified"}, http.StatusBadRequest)
		return
	}
	_, err = strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "ID must be a number"}, http.StatusBadRequest)
		return
	}

	task, err = checkTask(task)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	t, err := service.GetById(id)

	if t.ID == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
		return
	}

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	err = service.Update(task)
	if err != nil {
		log.Println("Update err=", err.Error())
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, emptyResponse, http.StatusOK)
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task ID not specified"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "Invalid task ID"}, http.StatusBadRequest)
		return
	}

	//check if the task exists
	t, err := service.GetById(id)
	if t.ID == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
		return
	}
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	err = service.Done(t)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, emptyResponse, http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task ID not specified"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: "Invalid task ID"}, http.StatusBadRequest)
		return
	}

	//check if the task exists
	t, err := service.GetById(id)
	if t.ID == "" {
		makeHttpResponse(w, ErrorResponse{Error: "Task not found"}, http.StatusBadRequest)
		return
	}
	if err != nil {
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	err = service.Delete(id)
	if err != nil {
		log.Println("Delete err=", err.Error())
		makeHttpResponse(w, ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	makeHttpResponse(w, emptyResponse, http.StatusOK)
}

func checkTask(t store.Task) (store.Task, error) {
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

func makeHttpResponse(w http.ResponseWriter, response any, status int) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.WriteHeader(status)
		w.Write(jsonResponse)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}
