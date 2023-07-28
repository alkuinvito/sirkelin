package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sirkelin/backend/models"
	"strings"
)

type UserRepository struct {
}

type IUserRepository interface {
	Get(db *gorm.DB) ([]models.User, error)
	GetByID(db *gorm.DB, uid string) (*models.User, error)
	GetByKeyword(db *gorm.DB, keyword string) ([]models.User, error)
	Save(db *gorm.DB, user *models.User) error
	Update(db *gorm.DB, user *models.User) error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Get(db *gorm.DB) ([]models.User, error) {
	var result []models.User
	err := db.Select("id", "fullname", "picture").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *UserRepository) GetByID(db *gorm.DB, uid string) (*models.User, error) {
	var result models.User
	err := db.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *UserRepository) GetByKeyword(db *gorm.DB, keyword string) ([]models.User, error) {
	var result []models.User
	err := db.Select("id", "fullname", "picture").Where("UPPER(fullname) LIKE ?", "%"+strings.ToUpper(keyword)+"%").Limit(5).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *UserRepository) Save(db *gorm.DB, user *models.User) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}

func (repo *UserRepository) Update(db *gorm.DB, user *models.User) error {
	return db.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error
}
