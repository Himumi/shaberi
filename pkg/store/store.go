package store

import (
	"gorm.io/gorm"
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
