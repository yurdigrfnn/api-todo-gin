package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Name      string    `gorm:"type:varchar(300)" json:"name"`
	Complete  bool      `gorm:"default:false" json:"complete"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

func (t *Todo) BeforeSave(tx *gorm.DB) error {
	now := time.Now()
	if t.ID == 0 {
		t.CreatedAt = now
	}
	t.UpdatedAt = now
	return nil
}
