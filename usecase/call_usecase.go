package usecase

import "Golang-CRUD/domain"

type CallUsecase struct {
	Repo domain.CallRepository
}

// Hàm constructor khởi tạo usecase
func NewCallUsecase(repo domain.CallRepository) *CallUsecase {
	return &CallUsecase{
		Repo: repo,
	}
}

// Create
func (u *CallUsecase) Create(call *domain.CallLog) error {
	return u.Repo.Create(call)
}

// Update
func (u *CallUsecase) Update(id uint, call *domain.CallLog) error {
	return u.Repo.Update(id, call)
}

// Delete
func (u *CallUsecase) Delete(id uint) error {
	return u.Repo.Delete(id)
}

// GetByID
func (u *CallUsecase) GetByID(id uint) (*domain.CallLog, error) {
	return u.Repo.GetByID(id)
}

// List
func (u *CallUsecase) List(filter domain.CallFilter) ([]domain.CallLog, error) {
	return u.Repo.List(filter)
}
