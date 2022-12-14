package statistics

import (
	"database/sql"
	"errors"
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

func (s *Service) GetCountAcceptedTask() (count int, err error) {
	count, err = s.db.GetCountAcceptedTask()
	if err != nil {
		return 0, fmt.Errorf("failed to get count accepted: %w", err)
	}
	return
}

func (s *Service) GetCountDeclinedTask() (count int, err error) {
	count, err = s.db.GetCountDeclinedTask()
	if err != nil {
		return 0, fmt.Errorf("failed to get count declined: %w", err)
	}
	return
}

func (s *Service) GetSumReaction(objectId uint32) (count int, err error) {
	count, err = s.db.GetSumReaction(objectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get count declined: %w", err)
	}
	return
}
