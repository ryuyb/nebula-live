package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"nebula-live/ent"
	"nebula-live/internal/infrastructure/config"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

// NewEntClient 创建Ent客户端
func NewEntClient(cfg *config.Config, logger *zap.Logger) (*ent.Client, error) {
	var db *sql.DB
	var dbDialect string
	var err error

	switch cfg.Database.Driver {
	case "sqlite":
		dbDialect = dialect.SQLite
		dsn := cfg.Database.Database

		// 如果不是内存数据库，确保目录存在
		if dsn != ":memory:" && dsn != "" {
			dir := filepath.Dir(dsn)
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return nil, fmt.Errorf("failed to create database directory: %w", err)
				}
			}
		}

		// 添加SQLite参数以启用外键约束和其他优化
		if dsn != ":memory:" {
			dsn += "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
		} else {
			dsn += "?_pragma=foreign_keys(1)"
		}

		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
		}

		logger.Info("SQLite database connection established successfully",
			zap.String("driver", cfg.Database.Driver),
			zap.String("database", cfg.Database.Database),
		)

	case "postgres", "postgresql":
		dbDialect = dialect.Postgres
		// 构建PostgreSQL连接字符串
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Database,
			cfg.Database.SSLMode,
		)

		// 使用pgx驱动打开数据库连接
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres connection: %w", err)
		}

		logger.Info("PostgreSQL database connection established successfully",
			zap.String("driver", cfg.Database.Driver),
			zap.String("host", cfg.Database.Host),
			zap.Int("port", cfg.Database.Port),
			zap.String("database", cfg.Database.Database),
		)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// 配置连接池（SQLite不需要连接池配置）
	if cfg.Database.Driver != "sqlite" {
		db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建Ent客户端
	drv := entsql.OpenDB(dbDialect, db)
	client := ent.NewClient(ent.Driver(drv))

	return client, nil
}

// RunMigrations 运行数据库迁移
func RunMigrations(ctx context.Context, client *ent.Client, logger *zap.Logger) error {
	logger.Info("Running database migrations")

	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// CloseEntClient 关闭Ent客户端
func CloseEntClient(client *ent.Client, logger *zap.Logger) error {
	logger.Info("Closing database connection")
	return client.Close()
}
