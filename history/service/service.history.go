package service

import (
	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/history"
)

// Service is a type that defines history service
type Service struct {
	historyRepo history.IHistoryRepository
}

// NewHistoryService is a function that returns a new history service
func NewHistoryService(historyRepository history.IHistoryRepository) history.IService {
	return &Service{historyRepo: historyRepository}
}

// AddHistory is a method that adds a new user history to the system
func (service *Service) AddHistory(newOPHistory *entity.UserHistory) error {
	return service.historyRepo.Create(newOPHistory)
}

// SearchHistories is a method that search and returns a set of user's histories that matchs the identifier value
func (service *Service) SearchHistories(key, orderBy string, methods []string, pageNum int, columns ...string) []*entity.UserHistory {
	return service.historyRepo.Search(key, orderBy, methods, pageNum, columns...)
}

// AllUserHistories is a method that returns all the histories that is related to a certain userID
func (service *Service) AllUserHistories(userID string) []*entity.UserHistory {
	return service.historyRepo.All(userID)
}

// FindHistory is a method that finds a certain history from the system using the identifer
func (service *Service) FindHistory(identifier int) (*entity.UserHistory, error) {
	return service.historyRepo.Find(identifier)
}
