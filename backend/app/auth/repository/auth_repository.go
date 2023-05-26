package repository

import (
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
	Email     string  `gorm:"uniqueIndex;not null"`
	Rooms     []*Room `gorm:"many2many:user_rooms"`
	CreatedAt time.Time
}

type AuthRepository struct {
	db *gorm.DB
}

type IAuthRepository interface {
	AuthenticateByIDToken(id, fullname, picture, email string) error
	GetUsers(uid string) ([]User, error)
	GetUsersByName(fullname string) ([]User, error)
	GetUserByID(uid string) (*User, error)
}

func Init() *AuthRepository {
	return &AuthRepository{db: initializers.DB}
}

func (repo *AuthRepository) AuthenticateByIDToken(id, fullname, picture, email string) error {
	user := &User{
		ID:       id,
		Fullname: fullname,
		Picture:  picture,
		Email:    email,
	}
	return repo.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}

func (repo *AuthRepository) GetUsers(uid string) ([]User, error) {
	var result []User
	err := repo.db.Select("id", "fullname", "picture").Not("id = ?", uid).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetUsersByName(fullname string) ([]User, error) {
	var result []User
	err := repo.db.Select("id", "fullname", "picture").Where("UPPER(fullname) LIKE ?", "%"+strings.ToUpper(fullname)+"%").Limit(5).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AuthRepository) GetUserByID(uid string) (*User, error) {
	var result User
	err := repo.db.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
