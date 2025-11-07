package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed request: %w", err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cannot get last ID: %w", err)
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {

	db := GetDB()
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`

	rows, err := db.Query(query, limit)
	if err != nil {
		log.Printf("failed request: %v", err)
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Printf("scar error: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("iteration error: %v", err)
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

var err error

func GetTask(id int) (*Task, error) {
	db := GetDB()
	err = nil
	if err != nil {
		log.Printf("db connection error: %v", err)
		return nil, err
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	task := &Task{}

	err = db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("task id=%d not found", id)
			return nil, fmt.Errorf("task id=%d not found", id)
		}
		log.Printf("failed request: %v", err)
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	db := GetDB()
	if err != nil {
		log.Printf("db connection error: %v", err)
		return err
	}
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("failed request: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows count error: %v", err)
		return err
	}

	if count == 0 {
		err := fmt.Errorf("task id=%s not found", task.ID)
		log.Println(err)
		return err
	}

	return nil
}

func DeleteTask(id string) error {
	db := GetDB()
	if err != nil {
		log.Printf("db connection error: %v", err)
		return err
	}

	query := "DELETE FROM scheduler WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		log.Printf("task delete error %v", err)
		return err
	}
	return nil
}

func UpdateDate(next string, id string) error {
	db := GetDB()
	if err != nil {
		log.Printf("db connection error: %v", err)
		return err
	}

	query := "UPDATE scheduler SET date = ? WHERE id = ?"

	res, err := db.Exec(query, next, id)
	if err != nil {
		log.Printf("failed request: %v", err)
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows count error: %v", err)
		return err
	}

	if count == 0 {
		err := fmt.Errorf("task id=%s not found", id)
		log.Println(err)
		return err
	}

	return nil
}
