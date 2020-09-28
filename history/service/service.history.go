package service

import (
	"errors"
	"regexp"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/history"
	"github.com/Benyam-S/onepay/notifier"
)

// Service is a type that defines history service
type Service struct {
	historyRepo history.IHistoryRepository
	notifier    *notifier.Notifier
}

// NewHistoryService is a function that returns a new history service
func NewHistoryService(historyRepository history.IHistoryRepository, historyChangeNotifier *notifier.Notifier) history.IService {
	return &Service{historyRepo: historyRepository, notifier: historyChangeNotifier}
}

// AddHistory is a method that adds a new user history to the system
func (service *Service) AddHistory(newOPHistory *entity.UserHistory) error {

	err := service.historyRepo.Create(newOPHistory)
	if err != nil {
		return errors.New("unable to add new history")
	}

	/* ++++++++++++++ NOTIFYING CHANGE +++++++++++++++ */
	service.notifier.NotifyHistoryChange(newOPHistory)
	/* +++++++++++++++++++++++++++++++++++++++++++++++ */

	return nil
}

// SearchHistories is a method that search and returns a set of user's histories that matchs the identifier value
func (service *Service) SearchHistories(key, orderBy string, methods []string, pageNum int64, columns ...string) ([]*entity.UserHistory, int64) {

	empty, _ := regexp.MatchString(`^\s*$`, key)
	if empty {
		return []*entity.UserHistory{}, 0
	}

	return service.historyRepo.Search(key, orderBy, methods, pageNum, columns...)
}

// AllUserHistories is a method that returns all the histories that is related to a certain userID
func (service *Service) AllUserHistories(userID string) []*entity.UserHistory {
	return service.historyRepo.All(userID)
}

// FindHistory is a method that finds a certain history from the system using the identifer
func (service *Service) FindHistory(identifier int64) (*entity.UserHistory, error) {

	history, err := service.historyRepo.Find(identifier)
	if err != nil {
		return nil, errors.New("history not found")
	}
	return history, nil
}

// MarkUserHistoriesAsSeen is a method that marks a certain user's histories as seen
func (service *Service) MarkUserHistoriesAsSeen(userID string) error {

	err := service.historyRepo.MarkAsSeen(userID)
	if err != nil {
		return errors.New("unable to update histories")
	}
	return nil
}
