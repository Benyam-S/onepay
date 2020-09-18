package history

import "github.com/Benyam-S/onepay/entity"

// IService is an interface that defines all the service methods of a history struct
type IService interface {
	AddHistory(newOPHistory *entity.UserHistory) error
	SearchHistories(key, orderBy string, methods []string, pageNum int64, columns ...string) ([]*entity.UserHistory, int64)
	FindHistory(identifier int64) (*entity.UserHistory, error)
	AllUserHistories(userID string) []*entity.UserHistory
	MarkUserHistoriesAsSeen(userID string) error
}
