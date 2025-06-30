package usecase

import "call-api/domain"

type UserUsecase struct {
	Repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{Repo: repo}
}

func (u *UserUsecase) GetByID(id uint) (*domain.User, error) {
	return u.Repo.GetByID(id)
}

func (u *UserUsecase) Update(id uint, user *domain.User) error {
	return u.Repo.Update(id, user)
}

func (u *UserUsecase) Delete(id uint) error {
	return u.Repo.Delete(id)
}
