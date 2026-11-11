package user

import (
	"uuid"

	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;index;"`
	Username  string         `json:"username" gorm:"unique"`
	Email     string         `json:"email" gorm:"unique"`
	Hash      string         `json:"-"`
	CreatedAt int64          `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt int64          `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
