package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"finalProject/pkg/db"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

func dataCheck(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(formatDate)
	}

	t, err := time.Parse(formatDate, task.Date)
	if err != nil {
		log.Println("wrong data format")
		return fmt.Errorf("error in date: %w", err)
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Println("wrong repeat")
			return fmt.Errorf("error in repeat: %w", err)
		}

		if !t.After(now) {
			task.Date = next
		} else {
			if t.Before(now.Truncate(24 * time.Hour)) {
				task.Date = now.Format(formatDate)
			}
			return nil
		}
	}

	if !t.After(now) {
		task.Date = now.Format(formatDate)
	}
	return nil
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandlerId(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
	}
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("deserializing JSON error:", err)
		writeJson(w, map[string]string{"error": "deserializing JSON error"})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("title is empty")
		writeJson(w, map[string]string{"error": "title is empty"})
		return
	}

	err = dataCheck(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("data check error:", err)
		writeJson(w, map[string]string{"error": "data check error"})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("нadd task error:", err)
		writeJson(w, map[string]string{"error": "add task error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, map[string]any{"id": id})
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("deserializing JSON error:", err)
		writeJson(w, map[string]string{"error": "deserializing JSON error"})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("title is empty")
		writeJson(w, map[string]string{"error": "title is empty"})
		return
	}

	if task.ID == "0" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("id is empty")
		writeJson(w, map[string]string{"error": "id is empty"})
		return
	}

	err = dataCheck(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("data check error:", err)
		writeJson(w, map[string]string{"error": "data check error"})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("update task error:", err)
		writeJson(w, map[string]string{"error": "update task error"})
		return
	}

	writeJson(w, map[string]any{})
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTaskHandlerId(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")
	if idString == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("id cannot be empty")
		writeJson(w, map[string]string{"error": "id cannot be empty"})
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("incorrect id:", err)
		writeJson(w, map[string]string{"error": "incorrect id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("task not found:", err)
		writeJson(w, map[string]string{"error": "task not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJson(w, task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")
	if idString == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("id cannot be empty")
		writeJson(w, map[string]string{"error": "id cannot be empty"})
		return
	}

	err := db.DeleteTask(idString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("delete task error:", err)
		writeJson(w, map[string]string{"error": "delete task error"})
		return
	}

	writeJson(w, map[string]any{})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("id cannot be empty")
		writeJson(w, map[string]string{"error": "id cannot be empty"})
		return
	}

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("wrong format:", err)
		writeJson(w, map[string]string{"error": "wrong format"})
		return

	}

	task, err := db.GetTask(idInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("task not found:", err)
		writeJson(w, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(idString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("delete task error:", err)
			writeJson(w, map[string]string{"error": "delete task error"})
		}
		writeJson(w, map[string]any{})
		return
	}

	nextDay := time.Now().AddDate(0, 0, 1)

	nextDate, err := NextDate(nextDay, task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error NextDate:", err)
		writeJson(w, map[string]string{"error": "error NextDate"})
		return
	}

	err = db.UpdateDate(nextDate, idString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("update data error:", err)
		writeJson(w, map[string]string{"error": "update data error"})
		return
	}

	writeJson(w, map[string]any{})

}
