package pg

import (
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

type Analytics struct {
	db *sql.DB
}

func NewAnalytics(db *sql.DB) Analytics {
	return Analytics{db: db}
}

func (a Analytics) CreateTask(objectId uint32) (taskId int, err error) {
	row := a.db.QueryRow("INSERT INTO tasks_events (object_id, created_at, finished_at) VALUES ($1, $2, $3) RETURNING id",
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
	row := a.db.QueryRow("INSERT INTO letters_events (task_id, email, sent_at, finished_at) VALUES ($1, $2, $3, $4) RETURNING id",
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
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func (a Analytics) GetSumReaction(objectId uint32) (sum int, err error) {
	var taskId int
	taskId, err = a.GetTaskByObjectId(objectId)
	if err != nil || taskId == 0 {
		return
	}
	err = a.db.QueryRow("SELECT SUM(finished_at-sent_at) FROM letters_events l WHERE is_accepted = TRUE AND task_id = $1", taskId).
		Scan(&sum)
	return
}
