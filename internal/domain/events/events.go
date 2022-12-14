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

func (s *Service) CreateTask(objectId uint32) error {
	err := s.db.CreateTask(objectId)
	if err != nil {
		return fmt.Errorf("failed to CreateTask: %w", err)
	}
	return nil
}

func (s *Service) DeleteTask(objectId uint32) error {
	err := s.db.DeleteTask(objectId)
	if err != nil {
		return fmt.Errorf("failed to DeleteTask: %w", err)
	}
	return nil
}

func (s *Service) FinishTask(objectId uint32) error {
	err := s.db.FinishTask(objectId)
	if err != nil {
		return fmt.Errorf("failed to FinishTask: %w", err)
	}
	return nil
}

func (s *Service) CreateLetter(objectId uint32, email string) error {
	err := s.db.CreateLetter(objectId, email)
	if err != nil {
		return fmt.Errorf("failed to CreateLetter: %w", err)
	}
	return nil
}

func (s *Service) AcceptedLetter(objectId uint32, email string) error {
	err := s.db.AcceptedLetter(objectId, email)
	if err != nil {
		return fmt.Errorf("failed to AcceptedLetter: %w", err)
	}
	return nil
}

func (s *Service) DeclinedLetter(objectId uint32, email string) error {
	err := s.db.DeclinedLetter(objectId, email)
	if err != nil {
		return fmt.Errorf("failed to DeclinedLetter: %w", err)
	}
	return nil
}
