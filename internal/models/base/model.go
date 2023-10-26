package base

import (
	"time"
)

// Model base model.
type Model struct {
	Id        uint       `json:"id"         gorm:"column:id;primary_key;auto_increment"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at"                    sql:"index"`
}
