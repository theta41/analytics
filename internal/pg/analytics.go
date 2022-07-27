package pg

import (
	"database/sql"
	"time"
)

type Analytics struct {
	db *sql.DB
}

func NewAnalytics(db *sql.DB) Analytics {
	return Analytics{db: db}
}

func (a Analytics) CreateTask(objectId uint32) (taskId int, err error) {
	row := a.db.QueryRow("INSERT INTO tasks_events (object_id, created_at, finished_at, is_done) VALUES ($1, $2, $3, FALSE) RETURNING id",
		objectId, time.Now().Unix(), 0)
	if row.Err() != nil {
		return 0, row.Err()
	}
	err = row.Scan(&taskId)
	if err != nil {
		return 0, err
	}
	return
}

func (a Analytics) DeleteTask(objectId uint32) (err error) {
	_, err = a.db.Exec("DELETE FROM tasks_events WHERE object_id = $1", objectId)
	return err
}

func (a Analytics) FinishTask(objectId uint32) (err error) {
	_, err = a.db.Exec("UPDATE tasks_events SET finished_at = $1, is_done = TRUE WHERE object_id = $2",
		time.Now().Unix(), objectId)
	return err
}

func (a Analytics) CreateLetter(objectId uint32, email string) (letterId int, err error) {
	var taskId int
	taskId, err = a.GetTaskByObjectId(objectId)
	if err != nil || taskId == 0 {
		return
	}
	row := a.db.QueryRow("INSERT INTO letters_events (task_id, email, sent_at, finished_at, is_accepted) VALUES ($1, $2, $3, $4, FALSE) RETURNING id",
		taskId, email, time.Now().Unix(), 0)
	if row.Err() != nil {
		return 0, row.Err()
	}
	err = row.Scan(&letterId)
	if err != nil {
		return 0, err
	}
	return
}

func (a Analytics) AcceptedLetter(objectId uint32, email string) error {
	var taskId int
	taskId, err := a.GetTaskByObjectId(objectId)
	if err != nil || taskId == 0 {
		return err
	}
	_, err = a.db.Exec("UPDATE letters_events SET finished_at = $1, is_accepted = TRUE WHERE task_id = $2 AND email = $3",
		time.Now().Unix(), taskId, email)
	return err
}

func (a Analytics) DeclinedLetter(objectId uint32, email string) error {
	var taskId int
	taskId, err := a.GetTaskByObjectId(objectId)
	if err != nil || taskId == 0 {
		return err
	}
	_, err = a.db.Exec("UPDATE letters_events SET finished_at = $1, is_accepted = FALSE WHERE task_id = $2 AND email = $3",
		time.Now().Unix(), taskId, email)
	return err
}

func (a Analytics) GetCountAcceptedTask() (count int, err error) {
	err = a.db.QueryRow("SELECT COUNT(id) FROM tasks_events t WHERE is_done = TRUE").Scan(&count)
	return
}

func (t Analytics) GetCountDeclinedTask() (count int, err error) {
	err = t.db.QueryRow("SELECT COUNT(id) FROM tasks_events t WHERE is_done = FALSE").Scan(&count)
	return
}

func (a Analytics) GetTaskByObjectId(objectId uint32) (taskId int, err error) {
	err = a.db.QueryRow("SELECT id FROM tasks_events t WHERE object_id = $1", objectId).Scan(&taskId)
	return
}

func (a Analytics) GetSumReaction(objectId uint32) (sum int, err error) {
	var taskId int
	taskId, err = a.GetTaskByObjectId(objectId)
	if err != nil || taskId == 0 {
		return
	}
	err = a.db.QueryRow("SELECT COALESCE(SUM(finished_at-sent_at), 0) FROM letters_events l WHERE is_accepted = TRUE AND task_id = $1", taskId).
		Scan(&sum)
	return
}
