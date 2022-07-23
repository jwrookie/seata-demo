package dao

import (
	"time"

	"gorm.io/gorm"
)

type A struct {
	AId       uint64 `gorm:"primaryKey"`
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
	DeletedAt uint64 `json:"-"`
}

type ADao struct {
}

func (a *ADao) Create(db *gorm.DB) error {
	now := uint64(time.Now().UnixNano() / int64(uint64(time.Millisecond)/uint64(time.Nanosecond)))
	data := A{
		CreatedAt: now,
		UpdatedAt: now,
	}

	return db.Create(&data).Error
}
