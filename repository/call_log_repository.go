package repository

import (
	"call-api/model"
	"gorm.io/gorm"
)

type CallLogRepository struct {
	DB *gorm.DB
}

func NewCallLogRepository(db *gorm.DB) *CallLogRepository {
	return &CallLogRepository{DB: db}
}

func (r *CallLogRepository) Create(log *model.CallLog) error {
	return r.DB.Create(log).Error
}

func (r *CallLogRepository) Update(id uint, updated *model.CallLog) error {
	return r.DB.Model(&model.CallLog{}).Where("id = ?", id).Updates(updated).Error
}

func (r *CallLogRepository) Delete(id uint) error {
	return r.DB.Delete(&model.CallLog{}, id).Error
}

func (r *CallLogRepository) GetByID(id uint) (*model.CallLog, error) {
	var log model.CallLog
	err := r.DB.First(&log, id).Error
	return &log, err
}
