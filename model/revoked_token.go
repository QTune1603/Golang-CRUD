package model

type RevokedToken struct {
	Token     string `gorm:"primaryKey"`
	RevokedAt int64
}
