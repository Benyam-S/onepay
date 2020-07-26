package history

import "github.com/Benyam-S/onepay/entity"

// IHistoryRepository is an interface that defines all the repository methods of a user history struct
type IHistoryRepository interface {
	Create(newOPHistory *entity.UserHistory) error
	Find(identifier int64) (*entity.UserHistory, error)
	Search(key, orderBy string, methods []string, pageNum int64, columns ...string) []*entity.UserHistory
	All(identifier string) []*entity.UserHistory
	Update(opHistory *entity.UserHistory) error
	Delete(identifier int64) (*entity.UserHistory, error)
}
