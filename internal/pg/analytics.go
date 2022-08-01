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

func (a Analytics) CreateTask(objectId uint32) error {
	row := a.db.QueryRow(
		"INSERT INTO tasks_events (object_id, created_at, finished_at, is_done) VALUES ($1, $2, $3, FALSE) RETURNING id",
		objectId,
		time.Now().Unix(),
		0,
	)

	if row.Err() != nil {
		return row.Err()
	}

	// var taskEventId int
	// err := row.Scan(&taskEventId)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a Analytics) DeleteTask(objectId uint32) error {
	_, err := a.db.Exec("DELETE FROM tasks_events WHERE object_id = $1", objectId)
	return err
}

func (a Analytics) FinishTask(objectId uint32) error {
	_, err := a.db.Exec("UPDATE tasks_events SET finished_at = $1, is_done = TRUE WHERE object_id = $2",
		time.Now().Unix(), objectId)
	return err
}

func (a Analytics) CreateLetter(objectId uint32, email string) error {
	taskEventId, err := a.getTaskByObjectId(objectId)
	if err != nil || taskEventId == 0 {
		return err
	}

	row := a.db.QueryRow(
		"INSERT INTO letters_events (task_event_id, email, sent_at, finished_at, is_accepted) VALUES ($1, $2, $3, $4, FALSE) RETURNING id",
		taskEventId,
		email,
		time.Now().Unix(),
		0,
	)

	if row.Err() != nil {
		return row.Err()
	}

	// var letterId int
	// err = row.Scan(&letterId)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (a Analytics) AcceptedLetter(objectId uint32, email string) error {
	taskEventId, err := a.getTaskByObjectId(objectId)
	if err != nil || taskEventId == 0 {
		return err
	}
	_, err = a.db.Exec(
		"UPDATE letters_events SET finished_at = $1, is_accepted = TRUE WHERE task_event_id = $2 AND email = $3",
		time.Now().Unix(),
		taskEventId,
		email,
	)
	return err
}

func (a Analytics) DeclinedLetter(objectId uint32, email string) error {
	taskEventId, err := a.getTaskByObjectId(objectId)
	if err != nil || taskEventId == 0 {
		return err
	}
	_, err = a.db.Exec(
		"UPDATE letters_events SET finished_at = $1, is_accepted = FALSE WHERE task_event_id = $2 AND email = $3",
		time.Now().Unix(),
		taskEventId,
		email,
	)
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

func (a Analytics) getTaskByObjectId(objectId uint32) (taskEventId int, err error) {
	err = a.db.QueryRow(
		"SELECT id FROM tasks_events t WHERE object_id = $1",
		objectId,
	).Scan(&taskEventId)
	return
}

func (a Analytics) GetSumReaction(objectId uint32) (sum int, err error) {
	var taskEventId int
	taskEventId, err = a.getTaskByObjectId(objectId)
	if err != nil || taskEventId == 0 {
		return
	}
	err = a.db.QueryRow(
		"SELECT COALESCE(SUM(finished_at-sent_at), 0) FROM letters_events l WHERE is_accepted = TRUE AND task_event_id = $1",
		taskEventId,
	).Scan(&sum)
	return
}
