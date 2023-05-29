package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
	"strings"
)

type AuthRepository struct {
	db *gorm.DB
}

type IAuthRepository interface {
	Get() ([]models.User, error)
	GetByID(uid string) (*models.User, error)
	GetByKeyword(keyword string) ([]models.User, error)
	GetExcept(uid string) ([]models.User, error)
	Save(user *models.User) error
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{db: initializers.DB}
}

func (repo *AuthRepository) Get() ([]models.User, error) {
	var result []models.User
	err := repo.db.Select("id", "fullname", "picture").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetByID(uid string) (*models.User, error) {
	var result models.User
	err := repo.db.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *AuthRepository) GetByKeyword(keyword string) ([]models.User, error) {
	var result []models.User
	err := repo.db.Select("id", "fullname", "picture").Where("UPPER(fullname) LIKE ?", "%"+strings.ToUpper(keyword)+"%").Limit(5).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetExcept(uid string) ([]models.User, error) {
	var result []models.User
	err := repo.db.Select("id", "fullname", "picture").Not("id = ?", uid).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) Save(user *models.User) error {
	return repo.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}
