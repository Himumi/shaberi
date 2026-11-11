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

func TestMigrateAndDropTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), SilentDB)
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	storeDB := New(db)

	// Test AutoMigrate
	err = storeDB.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	if !db.Migrator().HasTable(&TestModel{}) {
		t.Error("expected table TestModel to exist after migration")
	}

	// Test DropTable
	err = storeDB.DropTable(&TestModel{})
	if err != nil {
		t.Fatalf("failed to drop table: %v", err)
	}

	if db.Migrator().HasTable(&TestModel{}) {
		t.Error("expected table TestModel to be dropped")
	}
}
