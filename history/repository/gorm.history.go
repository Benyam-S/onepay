package repository

import (
	"fmt"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/history"
	"github.com/jinzhu/gorm"
)

// HistoryRepository is a type that defines a user history repository
type HistoryRepository struct {
	conn *gorm.DB
}

// NewHistoryRepository is a function that returns a new user history repository
func NewHistoryRepository(connection *gorm.DB) history.IHistoryRepository {
	return &HistoryRepository{conn: connection}
}

// Create is a method that adds a new user history to the database
func (repo *HistoryRepository) Create(newOPHistory *entity.UserHistory) error {

	err := repo.conn.Create(newOPHistory).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain user history from the database using an identifier.
// In Find() id is only used as a key
func (repo *HistoryRepository) Find(identifier int) (*entity.UserHistory, error) {
	opHistory := new(entity.UserHistory)
	err := repo.conn.Model(opHistory).
		Where("id = ?", identifier).First(opHistory).Error

	if err != nil {
		return nil, err
	}
	return opHistory, nil
}

// Search is a method that search and returns a set of user histories from the database using an identifier.
func (repo *HistoryRepository) Search(key, orderBy string, methods []string, pageNum int64, columns ...string) []*entity.UserHistory {

	var opHistories []*entity.UserHistory
	var whereStmt1 []string
	var whereStmt2 []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt1 = append(whereStmt1, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	for _, method := range methods {
		whereStmt2 = append(whereStmt2, fmt.Sprintf(" method = ? "))
		sqlValues = append(sqlValues, method)
	}

	sqlValues = append(sqlValues, orderBy)
	sqlValues = append(sqlValues, pageNum*20)

	repo.conn.Raw("SELECT * FROM user_history WHERE ("+strings.Join(whereStmt1, "||")+") && ("+strings.Join(whereStmt2, "||")+") ORDER BY ? DESC LIMIT ?, 20", sqlValues...).Scan(&opHistories)

	return opHistories
}

// All is a method that returns a set of user histories that is related to the key identifier
// In All() sender_id and receiver_id are used as an identifier
func (repo *HistoryRepository) All(identifier string) []*entity.UserHistory {

	var opHistories []*entity.UserHistory
	err := repo.conn.Model(entity.UserHistory{}).
		Where("sender_id = ? || receiver_id = ?", identifier, identifier).Find(&opHistories).Error

	if err != nil {
		return []*entity.UserHistory{}
	}
	return opHistories
}

// Update is a method that updates a certain user history value in the database
func (repo *HistoryRepository) Update(opHistory *entity.UserHistory) error {

	prevOPHistory := new(entity.UserHistory)
	err := repo.conn.Model(prevOPHistory).Where("id = ?", opHistory.ID).First(prevOPHistory).Error

	if err != nil {
		return err
	}

	err = repo.conn.Save(opHistory).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user history from the database using an identifier.
// In Delete() id is only used as a key
func (repo *HistoryRepository) Delete(identifier int) (*entity.UserHistory, error) {
	opHistory := new(entity.UserHistory)
	err := repo.conn.Model(opHistory).Where("id = ?", identifier).First(opHistory).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opHistory)
	return opHistory, nil
}
