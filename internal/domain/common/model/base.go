package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntityModel interface {
	TableName() string
	TracerName() string
}

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"primarykey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if bm.ID == uuid.Nil {
		bm.ID, _ = uuid.NewV7()
	}

	return nil
}
