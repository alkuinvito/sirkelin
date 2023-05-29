package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sirkelin/backend/models"
	"strings"
)

type AuthRepository struct {
}

type IAuthRepository interface {
	Get(db *gorm.DB) ([]models.User, error)
	GetByID(db *gorm.DB, uid string) (*models.User, error)
	GetByKeyword(db *gorm.DB, keyword string) ([]models.User, error)
	GetExcept(db *gorm.DB, uid string) ([]models.User, error)
	Save(db *gorm.DB, user *models.User) error
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (repo *AuthRepository) Get(db *gorm.DB) ([]models.User, error) {
	var result []models.User
	err := db.Select("id", "fullname", "picture").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetByID(db *gorm.DB, uid string) (*models.User, error) {
	var result models.User
	err := db.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *AuthRepository) GetByKeyword(db *gorm.DB, keyword string) ([]models.User, error) {
	var result []models.User
	err := db.Select("id", "fullname", "picture").Where("UPPER(fullname) LIKE ?", "%"+strings.ToUpper(keyword)+"%").Limit(5).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetExcept(db *gorm.DB, uid string) ([]models.User, error) {
	var result []models.User
	err := db.Select("id", "fullname", "picture").Not("id = ?", uid).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) Save(db *gorm.DB, user *models.User) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}
