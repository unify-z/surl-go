package user

import (
	"github.com/unify-z/go-surl/internal/repository/surl"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *UserModel) error {
	return u.db.Create(user).Error
}

func (u *UserRepo) GetUserByUsername(username string) (*UserModel, error) {
	var user UserModel
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetUserByID(id uint) (*UserModel, error) {
	var user UserModel
	err := u.db.Preload("URLs").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) UpdateUser(user *UserModel) error {
	return u.db.Save(user).Error
}

func (u *UserRepo) DeleteUser(id uint) error {
	return u.db.Delete(&UserModel{}, id).Error
}

func (u *UserRepo) ListUsersWithoutSURLs() ([]UserModel, error) {
	var users []UserModel
	err := u.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserRepo) GetUserSURLS(id uint) ([]surl.ShortURL, error) {
	var user UserModel
	err := u.db.Preload("URLs").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return user.URLs, nil

}

func (u *UserRepo) CountUsers() (int64, error) {
	var count int64
	err := u.db.Model(&UserModel{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
