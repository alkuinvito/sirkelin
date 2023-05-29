package repository

import (
	"github.com/alkuinvito/sirkelin/app/room/repository"
	"github.com/alkuinvito/sirkelin/initializers"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID        string
	Username  string
	Fullname  string
	Picture   string
	Email     string             `gorm:"uniqueIndex;not null"`
	Rooms     []*repository.Room `gorm:"many2many:user_rooms"`
	CreatedAt time.Time
}

type AuthRepository struct {
	db *gorm.DB
}

type IAuthRepository interface {
	Get() ([]User, error)
	GetByID(uid string) (*User, error)
	GetByKeyword(keyword string) ([]User, error)
	GetExcept(uid string) ([]User, error)
	Save(user *User) error
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{db: initializers.DB}
}

func (repo *AuthRepository) Get() ([]User, error) {
	var result []User
	err := repo.db.Select("id", "fullname", "picture").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetByID(uid string) (*User, error) {
	var result User
	err := repo.db.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *AuthRepository) GetByKeyword(keyword string) ([]User, error) {
	var result []User
	err := repo.db.Select("id", "fullname", "picture").Where("UPPER(fullname) LIKE ?", "%"+strings.ToUpper(keyword)+"%").Limit(5).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetExcept(uid string) ([]User, error) {
	var result []User
	err := repo.db.Select("id", "fullname", "picture").Not("id = ?", uid).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) Save(user *User) error {
	return repo.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}
