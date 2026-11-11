package user_test

import (
	"testing"
	"uuid"

	"github.com/Himumi/shaberi/internal/user"
	"github.com/Himumi/shaberi/pkg/phash"
	"github.com/Himumi/shaberi/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func TestUserBeforeUpdate(t *testing.T) {
	db := SetupTestDB(t)

	t.Run("Fails Update data is empty (map[string]any)", func(t *testing.T) {
		data := &user.User{}
		tx := db.Model(data).Updates(user.User{Username: "newname"})

		err := data.BeforeUpdate(tx)
		if err != nil {
			assert.ErrorIs(t, err, user.ErrEmptyUpdateData)
		}
	})

	t.Run("Fails Updating forbidden fields", func(t *testing.T) {
		cases := []struct {
			data map[string]any
		}{
			{data: map[string]any{
				user.FieldID: uuid.NewV4(),
			}},
			{data: map[string]any{
				user.FieldUsername: "newusername",
			}},
			{data: map[string]any{
				user.FieldEmail: "newemail@email.com",
			}},
		}

		for _, c := range cases {
			data := &user.User{}
			tx := db.Model(data).Updates(c.data)

			err := data.BeforeUpdate(tx)
			if err != nil {
				assert.ErrorIs(t, err, user.ErrCannotModifyData)
			}
		}
	})

	t.Run("Fails updating invalid password type (string)", func(t *testing.T) {
		data := &user.User{}
		tx := db.Model(data).Updates(map[string]any{
			user.FieldHash: 123456,
		})

		err := data.BeforeUpdate(tx)
		if err != nil {
			assert.ErrorIs(t, err, user.ErrInvalidType)
		}
	})

	t.Run("Succeeds hashing updated password", func(t *testing.T) {
		data := &user.User{}
		rawPassword := "NewPassword123!"

		tx := db.Model(data).Updates(map[string]any{
			user.FieldHash: rawPassword,
		})

		err := data.BeforeUpdate(tx)
		assert.NoError(t, err)

		updatedValues := tx.Statement.Dest.(map[string]any)
		hashed, ok := updatedValues[user.FieldHash].(string)

		require.True(t, ok)
		assert.NotEqual(t, rawPassword, hashed)
		assert.NotEmpty(t, hashed)

		err = phash.Verify(hashed, rawPassword)
		assert.NoError(t, err)
	})
}

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), store.SilentDB)
	require.NoError(t, err)

	err = db.AutoMigrate(&user.User{}, &user.Profile{})
	require.NoError(t, err)

	return db
}

func TestProfileBeforeCreate(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		profile := &user.Profile{
			Name:   "My name",
			Status: "Active",
			UserID: uuid.NewV4(),
		}

		err := profile.BeforeCreate(nil)
		require.NoError(t, err)
	})

	t.Run("Fails when User ID is empty", func(t *testing.T) {
		profile := &user.Profile{
			Name:   "My name",
			Status: "Active",
			UserID: uuid.Nil(),
		}

		err := profile.BeforeCreate(nil)
		if err != nil {
			assert.ErrorIs(t, err, user.ErrRequiredField)
		}
	})
}

func TestProfileBeforeUpdate(t *testing.T) {
	db := SetupTestDB(t)

	t.Run("Succeeds", func(t *testing.T) {
		profile := &user.Profile{}
		tx := db.Model(profile).Updates(map[string]any{
			user.FieldName:   "New profile name",
			user.FieldStatus: "Active",
		})

		err := profile.BeforeUpdate(tx)
		assert.NoError(t, err)
	})

	t.Run("Fails Updating forbiddin fields", func(t *testing.T) {
		cases := []map[string]any{
			{
				user.FieldID: "new id",
			},
			{
				user.FieldUserID: uuid.NewV4(),
			},
		}

		for _, c := range cases {
			profile := &user.Profile{}
			tx := db.Model(profile).Updates(c)

			err := profile.BeforeUpdate(tx)
			if err != nil {
				assert.ErrorIs(t, err, user.ErrCannotModifyData)
			}
		}
	})
}
