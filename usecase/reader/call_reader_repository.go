package reader

import "Golang-CRUD/domain"

// CallReaderRepository chỉ phục vụ filter metadata
type CallReaderRepository interface {
	GetWithMetadataField(filter domain.CallFilter, metaField string) ([]domain.CallLog, error)
}
