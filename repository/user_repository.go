package repository

import (
	"call-api/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	Password  string
	CreatedAt int64
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *domain.User) error {
	model := UserModel{
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
	return r.db.Create(&model).Error
}

func (r *userRepository) Update(id uint, user *domain.User) error {
	updateData := map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
	}
	return r.db.Model(&UserModel{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&UserModel{}, id).Error
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var model UserModel
	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        model.ID,
		Username:  model.Username,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var model UserModel
	err := r.db.Where("username = ?", username).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        model.ID,
		Username:  model.Username,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}, nil
}
