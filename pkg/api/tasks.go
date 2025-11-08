package api

import (
	"log"
	"net/http"

	"finalProject/pkg/db"
)

const limit = 50

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.Tasks(limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getting tasks error:", err)
		writeJson(w, map[string]string{"error": "getting tasks error"})
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	w.WriteHeader(http.StatusOK)
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
