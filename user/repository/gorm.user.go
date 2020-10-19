package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Benyam-S/onepay/entity"
	"github.com/Benyam-S/onepay/tools"
	"github.com/Benyam-S/onepay/user"
	"github.com/jinzhu/gorm"
)

// UserRepository is a type that defines a user repository
type UserRepository struct {
	conn *gorm.DB
}

// NewUserRepository is a function that returns a new user repository
func NewUserRepository(connection *gorm.DB) user.IUserRepository {
	return &UserRepository{conn: connection}
}

// Create is a method that adds a new user to the database
func (repo *UserRepository) Create(newOPUser *entity.User) error {
	totalNumOfUsers := repo.CountUsers()
	baseID := 1000101101010
	newOPUser.UserID = fmt.Sprintf("OP-%d", baseID+totalNumOfUsers)

	for !repo.IsUnique("user_id", newOPUser.UserID) {
		totalNumOfUsers++
		newOPUser.UserID = fmt.Sprintf("OP-%d", baseID+totalNumOfUsers)
	}

	err := repo.conn.Create(newOPUser).Error
	if err != nil {
		return err
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	repo.conn.Exec("UPDATE extras SET total_users_count = ?", totalNumOfUsers+1)
	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */

	return nil
}

// Find is a method that finds a certain user from the database using an identifier,
// also Find() uses user_id and email as a key for selection
func (repo *UserRepository) Find(identifier string) (*entity.User, error) {
	opUser := new(entity.User)

	err := repo.conn.Model(opUser).Where("user_id = ? || email = ?",
		identifier, identifier).First(opUser).Error

	if err != nil {
		return nil, err
	}
	return opUser, nil
}

// FindAlsoWPhone is a method that finds a certain user from the database using an identifier,
// also FindAlsoWPhone() uses user_id, email and phone_number as a key for selection
func (repo *UserRepository) FindAlsoWPhone(identifier, phoneNumber string) (*entity.User, error) {
	phoneNumber = `^` + tools.EscapeRegexpForDatabase(phoneNumber) + `(\\[[a-zA-Z]{2}])?$`

	opUser := new(entity.User)
	err := repo.conn.Model(opUser).Where("user_id = ? || email = ? || phone_number REGEXP '"+
		phoneNumber+"'", identifier, identifier).First(opUser).Error

	if err != nil {
		return nil, err
	}
	return opUser, nil
}

// Update is a method that updates a certain user value in the database
func (repo *UserRepository) Update(opUser *entity.User) error {

	prevOPUser := new(entity.User)
	err := repo.conn.Model(prevOPUser).Where("user_id = ?", opUser.UserID).First(prevOPUser).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */
	if opUser.ProfilePic == "" {
		opUser.ProfilePic = prevOPUser.ProfilePic
	}
	opUser.CreatedAt = prevOPUser.CreatedAt
	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(opUser).Error
	if err != nil {
		return err
	}
	return nil
}

// Search is a method that search and returns a set of users from the database using an identifier.
func (repo *UserRepository) Search(key string, pageNum int64, columns ...string) []*entity.User {

	var opUsers []*entity.User
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		// modifing the key so that it can match the database phone number values
		if column == "phone_number" {
			modifiedKey := key
			splitedKey := strings.Split(key, "")
			if splitedKey[0] == "0" && len(splitedKey) == 10 {
				modifiedKey = "+251" + strings.Join(splitedKey[1:], "")
			}
			modifiedKey = `^` + tools.EscapeRegexpForDatabase(modifiedKey) + `(\\[[a-zA-Z]{2}])?$`
			whereStmt = append(whereStmt, fmt.Sprintf(" %s REGEXP '"+modifiedKey+"' ", column))
			continue
		}
		whereStmt = append(whereStmt, fmt.Sprintf(" %s = ? ", column))
		sqlValues = append(sqlValues, key)
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM users WHERE ("+strings.Join(whereStmt, "||")+") ORDER BY first_name ASC LIMIT ?, 30", sqlValues...).Scan(&opUsers)

	return opUsers
}

// SearchWRegx is a method that searchs and returns set of users limited to the key identifier and page number using regular experssions
func (repo *UserRepository) SearchWRegx(key string, pageNum int64, columns ...string) []*entity.User {
	var opUsers []*entity.User
	var whereStmt []string
	var sqlValues []interface{}

	for _, column := range columns {
		whereStmt = append(whereStmt, fmt.Sprintf(" %s REGEXP ? ", column))
		sqlValues = append(sqlValues, "^"+regexp.QuoteMeta(key))
	}

	sqlValues = append(sqlValues, pageNum*30)
	repo.conn.Raw("SELECT * FROM users WHERE "+strings.Join(whereStmt, "||")+" ORDER BY first_name ASC LIMIT ?, 30", sqlValues...).Scan(&opUsers)

	return opUsers
}

// All is a method that returns all the users from the database limited with the pageNum
func (repo *UserRepository) All(pageNum int64) []*entity.User {

	var opUsers []*entity.User
	limit := pageNum * 30

	repo.conn.Raw("SELECT * FROM users ORDER BY first_name ASC LIMIT ?, 30", limit).Scan(&opUsers)
	return opUsers
}

// UpdateValue is a method that updates a certain user's single column value in the database
func (repo *UserRepository) UpdateValue(opUser *entity.User, columnName string, columnValue interface{}) error {

	prevOPUser := new(entity.User)
	err := repo.conn.Model(prevOPUser).Where("user_id = ?", opUser.UserID).First(prevOPUser).Error

	if err != nil {
		return err
	}

	err = repo.conn.Model(entity.User{}).Where("user_id = ?", opUser.UserID).Update(map[string]interface{}{columnName: columnValue}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete is a method that deletes a certain user from the database using an identifier.
// In Delete() user_id is only used as a key
func (repo *UserRepository) Delete(identifier string) (*entity.User, error) {
	opUser := new(entity.User)
	err := repo.conn.Model(opUser).Where("user_id = ?", identifier).First(opUser).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(opUser)
	return opUser, nil
}

// CountUsers is a method that counts the users in the database
func (repo *UserRepository) CountUsers() int {
	extrasColumn := make([]int, 0)
	repo.conn.Model(&entity.Extras{}).Pluck("total_users_count", &extrasColumn)
	return extrasColumn[0]
}

// IsUnique is a method that determines whether a certain column value is unique in the user table
func (repo *UserRepository) IsUnique(columnName string, columnValue interface{}) bool {
	var totalCount int
	repo.conn.Model(&entity.User{}).Where(columnName+"=?", columnValue).Count(&totalCount)
	return 0 >= totalCount
}

// IsUniqueRexp is a method that determines whether a certain column value pattern is unique in the user table
func (repo *UserRepository) IsUniqueRexp(columnName string, columnPattern string) bool {
	var totalCount int
	repo.conn.Raw("SELECT COUNT(*) FROM users WHERE " + columnName +
		" REGEXP '" + columnPattern + "'").Count(&totalCount)
	return 0 >= totalCount
}
