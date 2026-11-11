package user

import (
	"uuid"

	"github.com/Himumi/shaberi/pkg/phash"
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

// Possible Errors:
//   - ErrTooShort
//   - ErrTooLong
//   - bcrypt.GenerateFromPassword's Errors
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil() {
		u.ID = uuid.NewV4()
	}

	hashed, err := phash.Generate(u.Hash)
	if err != nil {
		return err
	}
	u.Hash = string(hashed)

	return nil
}
