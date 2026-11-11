package user

import (
	"errors"
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

var (
	ErrEmptyUpdateData  = errors.New("Update data is empty")
	ErrCannotModifyData = errors.New("Cannot to update this data")
	ErrInvalidType      = errors.New("Invalid type")
	FieldID             = "ID"
	FieldUsername       = "Username"
	FieldEmail          = "Email"
	FieldHash           = "Hash"
	ForbiddenFields     = []string{FieldID, FieldUsername, FieldEmail}
)

// Possible Errors:
//   - ErrEmptyUpdateData
//   - ErrCannotModifyData
//   - ErrInvalidType
//   - ErrTooShort
//   - ErrTooLong
//   - bcrypt.GenerateFromPassword's Errors
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	values, ok := tx.Statement.Dest.(map[string]any)
	if !ok {
		return ErrEmptyUpdateData
	}

	for _, field := range ForbiddenFields {
		if tx.Statement.Changed(field) {
			return ErrCannotModifyData
		}
	}

	if tx.Statement.Changed(FieldHash) {
		value, ok := values[FieldHash].(string)
		if !ok {
			return ErrInvalidType
		}

		hashed, err := phash.Generate(value)
		if err != nil {
			return err
		}

		tx.Statement.SetColumn(FieldHash, string(hashed))
	}

	return nil
}
