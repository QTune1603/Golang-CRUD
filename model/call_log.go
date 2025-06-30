package model

import (
	"gorm.io/datatypes"
)

type CallLog struct {
	ID uint `json:"id" gorm:"primaryKey"`
	PhoneNumber string `json:"phone_number"`
	Metadata datatypes.JSON `json:"metadata" gorm:"type:json"`
	CallResult string `json:"call_result"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	CallTime int64 `json:"call_time"`
	ResultTime int64 `json:"result_time"`
	PickupTime *int64 `json:"pickup_time"`
	HangupTime *int64 `json:"hangup_time"`
}