package repository

import (
	"fmt"
	"strings"

	"github.com/Benyam-S/onepay/deleted"
	"github.com/Benyam-S/onepay/entity"
	"github.com/jinzhu/gorm"
)

// FrozenClientRepository is a type that defines a repository for frozen api client
type FrozenClientRepository struct {
	conn *gorm.DB
}

// NewFrozenClientRepository is a function that returns a new frozen api client repository
func NewFrozenClientRepository(connection *gorm.DB) deleted.IFrozenClientRepository {
	return &FrozenClientRepository{conn: connection}
}

// Create is a method that adds a frozen api client to the database
func (repo *FrozenClientRepository) Create(frozenAPIClient *entity.FrozenClient) error {

	err := repo.conn.Create(frozenAPIClient).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is a method that finds a certain frozen api client from the database using an identifier,
// also Find() uses api_key as a key for selection
func (repo *FrozenClientRepository) Find(identifier string) (*entity.FrozenClient, error) {

	frozenAPIClient := new(entity.FrozenClient)
	err := repo.conn.Model(frozenAPIClient).
		Where("api_key = ? ", identifier).
		First(frozenAPIClient).Error

	if err != nil {
		return nil, err
	}
	return frozenAPIClient, nil
}

// Search is a method that search and returns a set of frozen api clients from the database using an identifier.
func (repo *FrozenClientRepository) Search(key string, pageNum int64, columns ...string) []*entity.FrozenClient {

	var frozenAPIClients []*entity.FrozenClient
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM frozen_clients WHERE ("+strings.Join(whereStmt, "||")+") ORDER BY api_key ASC LIMIT ?, 30", sqlValues...).Scan(&frozenAPIClients)

	return frozenAPIClients
}

// All is a method that returns all the api clients from the database limited with the pageNum
func (repo *FrozenClientRepository) All(pageNum int64) []*entity.FrozenClient {

	var frozenAPIClients []*entity.FrozenClient
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM frozen_clients ORDER BY api_key ASC LIMIT ?, 30", limit).Scan(&frozenAPIClients)
	return frozenAPIClients
}

// Update is a method that updates a certain frozen api client value in the database
func (repo *FrozenClientRepository) Update(frozenAPIClient *entity.FrozenClient) error {

	prevOPUser := new(entity.FrozenClient)
	err := repo.conn.Model(prevOPUser).Where("api_key = ?", frozenAPIClient.APIKey).First(prevOPUser).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	frozenAPIClient.CreatedAt = prevOPUser.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(frozenAPIClient).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain frozen api client from the database using an identifier.
// In Delete() api_key is only used as a key
func (repo *FrozenClientRepository) Delete(identifier string) (*entity.FrozenClient, error) {
	frozenAPIClient := new(entity.FrozenClient)
	err := repo.conn.Model(frozenAPIClient).Where("api_key = ?", identifier).First(frozenAPIClient).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(frozenAPIClient)
	return frozenAPIClient, nil
}
