package store

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/postgres"
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

var SilentDB = &gorm.Config{
	Logger: logger.Default.LogMode(logger.Silent),
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

var DefaultPoolConfig = PoolConfig{
	MaxOpenConns:    25, // Max postgres connections are 100
	MaxIdleConns:    10, // Keep 40% connections idle
	ConnMaxLifetime: 30 * time.Minute,
	ConnMaxIdleTime: 5 * time.Minute,
}

var (
	ErrDBConnection  = errors.New("Failed to connect to database")
	ErrDBInvalid     = errors.New("Invalid database")
	ErrDBUnreachable = errors.New("Unreachable database")
)

// Connect to the postgres
//
// Possible Errors:
//   - ErrDBConnection
//   - ErrDBInvalid
//   - ErrDBUnreachable
func ConnectPool(ctx context.Context, dsn string, poolCfg PoolConfig, opts ...gorm.Option) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), opts...)
	if err != nil {
		return nil, ErrDBConnection
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, ErrDBInvalid
	}

	if poolCfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(poolCfg.MaxOpenConns)
	}
	if poolCfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(poolCfg.MaxIdleConns)
	}
	if poolCfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(poolCfg.ConnMaxLifetime)
	}
	if poolCfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(poolCfg.ConnMaxIdleTime)
	}

	pingCtx, cancelCtx := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtx()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, ErrDBUnreachable
	}

	return db, nil
}
