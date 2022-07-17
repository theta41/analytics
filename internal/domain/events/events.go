package events

import (
	"fmt"

	"gitlab.com/g6834/team41/analytics/internal/repositories"
)

type Service struct {
	db repositories.Analytics
}

func New(db repositories.Analytics) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateTask(objectId uint32) (taskId int, err error) {
	taskId, err = s.db.CreateTask(objectId)
	if err != nil {
		return 0, fmt.Errorf("Failed to CreateTask: %w", err)
	}
	return
}

func (s *Service) FinishTask(objectId uint32) (err error) {
	err = s.db.FinishTask(objectId)
	if err != nil {
		return fmt.Errorf("Failed to FinishTask: %w", err)
	}
	return
}

func (s *Service) CreateLetter(objectId uint32, email string) (letterId int, err error) {
	letterId, err = s.db.CreateLetter(objectId, email)
	if err != nil {
		return 0, fmt.Errorf("Failed to CreateLetter: %w", err)
	}
	return
}

func (s *Service) AcceptedLetter(objectId uint32, email string) (err error) {
	err = s.db.AcceptedLetter(objectId, email)
	if err != nil {
		return fmt.Errorf("Failed to AcceptedLetter: %w", err)
	}
	return
}

func (s *Service) DeclinedLetter(objectId uint32, email string) (err error) {
	err = s.db.DeclinedLetter(objectId, email)
	if err != nil {
		return fmt.Errorf("Failed to DeclinedLetter: %w", err)
	}
	return
}
