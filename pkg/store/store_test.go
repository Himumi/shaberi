package store

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	gorm.Model
	Name string
}

func TestNewAndDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	if db == nil {
		t.Error("expected DB instance to not be nil")
	}

	storeDB := New(db)
	var _ Storer = storeDB
}
