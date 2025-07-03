// internal/infra/repository/call_reader_repository.go
package repository

import (
	"Golang-CRUD/domain"
	"Golang-CRUD/usecase/reader"
	"encoding/json"
	"gorm.io/gorm"
)

type callReaderRepository struct {
	db *gorm.DB
}

func NewCallReaderRepository(db *gorm.DB) reader.CallReaderRepository {
	return &callReaderRepository{db: db}
}

func (r *callReaderRepository) GetWithMetadataField(filter domain.CallFilter, metaField string) ([]domain.CallLog, error) {
	var models []CallModel

	query := r.db.Model(&CallModel{})

	if filter.PhoneNumber != "" {
		query = query.Where("phone_number LIKE ?", "%"+filter.PhoneNumber+"%")
	}
	if filter.StartAt != 0 && filter.EndAt != 0 {
		query = query.Where("created_at BETWEEN ? AND ?", filter.StartAt, filter.EndAt)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	var logs []domain.CallLog
	for _, m := range models {
		var meta map[string]interface{}
		_ = json.Unmarshal(m.Metadata, &meta)

		if metaField != "" {
			if _, ok := meta[metaField]; !ok {
				continue
			}
		}

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
