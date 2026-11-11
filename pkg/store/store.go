package store

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(db *gorm.DB) *Store {
	return &Store{db}
}

type Storer interface {
	DB() *gorm.DB
	Migrate(models ...any) error
	DropTable(tables ...any) error
}

type Store struct {
	db *gorm.DB
}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (s *Store) Migrate(models ...any) error {
	return s.DB().AutoMigrate(models...)
}

func (s *Store) DropTable(tables ...any) error {
	return s.DB().Migrator().DropTable(tables...)
}

func (s *Store) SetSilentMode() {
	s.db = s.DB().Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
