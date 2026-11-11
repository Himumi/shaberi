package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
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

	storeDB := New(db, nil)
	var _ Storer = storeDB
}

func TestMigrateAndDropTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), SilentDB)
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	storeDB := New(db, nil)

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

func TestSetSilentMode(t *testing.T) {
	t.Run("both database and redis client active", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to connect database: %v", err)
		}

		mr := miniredis.RunT(t)
		rdb := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		defer rdb.Close()

		store := New(db, rdb)
		store.SetSilentMode()

		if store.DB() == nil {
			t.Error("expected DB instance to remain valid")
		}
		if store.RDB() == nil {
			t.Error("expected RDB instance to remain valid")
		}
	})

	t.Run("only database, redis is nil", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to connect database: %v", err)
		}

		store := New(db, nil)
		store.SetSilentMode()

		if store.DB() == nil {
			t.Error("expected DB instance to remain valid")
		}
		if store.RDB() != nil {
			t.Error("expected RDB instance to be nil")
		}
	})

	t.Run("only redis, database is nil", func(t *testing.T) {
		mr := miniredis.RunT(t)
		rdb := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		defer rdb.Close()

		store := New(nil, rdb)
		store.SetSilentMode()

		if store.DB() != nil {
			t.Error("expected DB instance to be nil")
		}
		if store.RDB() == nil {
			t.Error("expected RDB instance to remain valid")
		}
	})
}

func TestConnectPool(t *testing.T) {
	t.Run("Invalid DSN", func(t *testing.T) {
		ctx := context.Background()
		invalidDSN := "host=127.0.0.1 port=1 user=test password=test dbname=test sslmode=disable"

		db, err := ConnectPool(ctx, invalidDSN, DefaultPoolConfig, SilentDB)
		if err == nil {
			t.Error("expected an error with an invalid DSN, got nil")
		}
		if db != nil {
			t.Error("expected db to be nil when connection fails")
		}
	})

	t.Run("Context canceled", func(t *testing.T) {
		ctx, cancelCtx := context.WithCancel(context.Background())
		cancelCtx() // Cancel immediately

		dsn := "host=192.0.2.1 port=5432 user=test password=test dbname=test sslmode=disable connect_timeout=1"

		db, err := ConnectPool(ctx, dsn, DefaultPoolConfig, SilentDB)
		if err == nil {
			t.Error("expected an error due to canceled context, got nil")
		}
		if db != nil {
			t.Error("expected db to be nil on context cancellation")
		}
	})

	t.Run("Success", func(t *testing.T) {
		dsn := os.Getenv("TEST_POSTGRES_DSN")
		if dsn == "" {
			t.Skip("Skipping postgres pool integration test: TEST_POSTGRES_DSN not set")
		}

		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelCtx()

		customConfig := PoolConfig{
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 10 * time.Minute,
			ConnMaxIdleTime: 2 * time.Minute,
		}

		db, err := ConnectPool(ctx, dsn, customConfig, SilentDB)
		if err != nil {
			t.Fatalf("failed to connect with pool config: %v", err)
		}
		if db == nil {
			t.Error("expected db instance to not be nil")
		}

		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get underlying sql.DB: %v", err)
		}

		if err := sqlDB.Ping(); err != nil {
			t.Errorf("expected ping to succeed: %v", err)
		}
	})
}

func TestConnectCache(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mr := miniredis.RunT(t)

		opt, err := redis.ParseURL("redis://" + mr.Addr())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		client, err := ConnectCache(t.Context(), opt)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer client.Close()

		if client == nil {
			t.Fatal("expected redis client to be non-nil")
		}
	})

	t.Run("Unreachable server", func(t *testing.T) {
		ctx, cancelCtx := context.WithCancel(t.Context())
		cancelCtx()

		opt, err := redis.ParseURL("redis://127.0.0.1:1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		client, err := ConnectCache(ctx, opt)
		if err == nil {
			t.Fatal("expected connection error, got nil")
		}
		if client != nil {
			client.Close()
			t.Fatal("expected nil client on connection failure")
		}
	})
}
