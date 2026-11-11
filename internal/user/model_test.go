package user_test

import (
	"testing"
	"uuid"

	"github.com/Himumi/shaberi/internal/user"
	"github.com/Himumi/shaberi/pkg/phash"
	"github.com/stretchr/testify/assert"
)

func TestUserBeforeCreate(t *testing.T) {
	t.Run("Generate ID", func(t *testing.T) {
		data := &user.User{
			Username: "testuser",
			Email:    "test@email.com",
			Hash:     "ValidPassword123!",
		}

		err := data.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil(), data.ID)
	})

	t.Run("Not generate ID if it exists", func(t *testing.T) {
		existingID := uuid.NewV4()
		data := &user.User{
			ID:       existingID,
			Username: "testuser",
			Email:    "test@email.com",
			Hash:     "ValidPassword123!",
		}

		err := data.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.Equal(t, existingID, data.ID)
	})

	t.Run("Succeeds hashing password", func(t *testing.T) {
		rawPassword := "SecurePassword123!"
		data := &user.User{
			Username: "testuser",
			Email:    "test@email.com",
			Hash:     rawPassword,
		}

		err := data.BeforeCreate(nil)
		assert.NoError(t, err)
		assert.NotEqual(t, rawPassword, data.Hash)
		assert.NotEmpty(t, data.Hash)
	})

	t.Run("Fails when password is short", func(t *testing.T) {
		data := &user.User{
			Username: "testuser",
			Email:    "test@email.com",
			Hash:     "short",
		}

		err := data.BeforeCreate(nil)
		if err != nil {
			assert.ErrorIs(t, err, phash.ErrTooShort)
		}
	})
}
