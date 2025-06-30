package repository

import (
	"call-api/domain"
	"encoding/json"
	"gorm.io/gorm"
)

type CallModel struct {
	gorm.Model
	PhoneNumber string
	Metadata    []byte
	CallResult  string
	CreatedAt   int64
	UpdatedAt   int64
	CallTime    int64
	ResultTime  int64
	PickupTime  *int64
	HangupTime  *int64
}

type callRepository struct {
	db *gorm.DB
}

func NewCallRepository(db *gorm.DB) domain.CallRepository {
	return &callRepository{db: db}
}

func (r *callRepository) Create(log *domain.CallLog) error {
	metaBytes, _ := json.Marshal(log.Metadata)

	model := CallModel{
		PhoneNumber: log.PhoneNumber,
		Metadata:    metaBytes,
		CallResult:  string(log.CallResult),
		CreatedAt:   log.CreatedAt,
		UpdatedAt:   log.UpdatedAt,
		CallTime:    log.CallTime,
		ResultTime:  log.ResultTime,
		PickupTime:  log.PickupTime,
		HangupTime:  log.HangupTime,
	}

	return r.db.Create(&model).Error
}

func (r *callRepository) Update(id uint, updated *domain.CallLog) error {
	metaBytes, _ := json.Marshal(updated.Metadata)

	updateData := map[string]interface{}{
		"phone_number": updated.PhoneNumber,
		"metadata":     metaBytes,
		"call_result":  string(updated.CallResult),
		"updated_at":   updated.UpdatedAt,
		"call_time":    updated.CallTime,
		"result_time":  updated.ResultTime,
		"pickup_time":  updated.PickupTime,
		"hangup_time":  updated.HangupTime,
	}

	return r.db.Model(&CallModel{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *callRepository) Delete(id uint) error {
	return r.db.Delete(&CallModel{}, id).Error
}

func (r *callRepository) GetByID(id uint) (*domain.CallLog, error) {
	var model CallModel
	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}

	var meta map[string]interface{}
	_ = json.Unmarshal(model.Metadata, &meta)

	return &domain.CallLog{
		ID:          model.ID,
		PhoneNumber: model.PhoneNumber,
		Metadata:    meta,
		CallResult:  domain.CallResult(model.CallResult),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CallTime:    model.CallTime,
		ResultTime:  model.ResultTime,
		PickupTime:  model.PickupTime,
		HangupTime:  model.HangupTime,
	}, nil
}

func (r *callRepository) List(filter domain.CallFilter) ([]domain.CallLog, error) {
	var models []CallModel
	query := r.db.Model(&CallModel{})

	if filter.PhoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+filter.PhoneNumber+"%")
	}
	if filter.StartAt != 0 && filter.EndAt != 0 {
		query = query.Where("created_at BETWEEN ? AND ?", filter.StartAt, filter.EndAt)
	}

	err := query.Find(&models).Error
	if err != nil {
		return nil, err
	}

	var logs []domain.CallLog
	for _, m := range models {
		var meta map[string]interface{}
		_ = json.Unmarshal(m.Metadata, &meta)

		logs = append(logs, domain.CallLog{
			ID:          m.ID,
			PhoneNumber: m.PhoneNumber,
			Metadata:    meta,
			CallResult:  domain.CallResult(m.CallResult),
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			CallTime:    m.CallTime,
			ResultTime:  m.ResultTime,
			PickupTime:  m.PickupTime,
			HangupTime:  m.HangupTime,
		})
	}

	return logs, nil
}
