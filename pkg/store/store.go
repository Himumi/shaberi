package store

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(db *gorm.DB, rdb *redis.Client) *Store {
	return &Store{db, rdb}
}

type Storer interface {
	DB() *gorm.DB
	Migrate(models ...any) error
	DropTable(tables ...any) error

	RDB() *redis.Client
}

type Store struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (s *Store) RDB() *redis.Client {
	return s.rdb
}

func (s *Store) Migrate(models ...any) error {
	return s.DB().AutoMigrate(models...)
}

func (s *Store) DropTable(tables ...any) error {
	return s.DB().Migrator().DropTable(tables...)
}

func (s *Store) Flush(ctx context.Context) error {
	return s.RDB().FlushAll(ctx).Err()
}

// For testing purposes
func (s *Store) SetSilentMode() {
	if s.DB() != nil {
		s.db = s.DB().Session(&gorm.Session{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	}
	if s.RDB() != nil {
		redis.SetLogger(silentLogger{})
	}
}

type silentLogger struct{}

func (silentLogger) Printf(ctx context.Context, format string, v ...any) {
	// Do nothing, effectively disabling logging
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

// Connect to the redis
//
// Possible Errors:
//   - ErrDBUnreachable
func ConnectCache(ctx context.Context, opt *redis.Options) (*redis.Client, error) {
	rdb := redis.NewClient(opt)

	ctx, cancelCtx := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtx()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, ErrDBUnreachable
	}

	return rdb, nil
}
