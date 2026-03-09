package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RedPanda13/task_manager/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func NewPostgres(ctx context.Context, cfg config.Config) (*Postgres, error) {
	dsn := buildDSN(cfg)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 buildLogger(cfg),
	})
	if err != nil {
		return nil, fmt.Errorf("open gorm connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("extract sql.DB from gorm: %w", err)
	}

	configurePool(sqlDB, cfg)

	if err := ping(ctx, sqlDB, cfg); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Postgres{
		Gorm: gormDB,
		SQL:  sqlDB,
	}, nil
}

func (p *Postgres) Close() error {
	if p == nil || p.SQL == nil {
		return nil
	}
	return p.SQL.Close()
}

func buildDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)
}

func configurePool(db *sql.DB, cfg config.Config) {
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)
}

func ping(parent context.Context, db *sql.DB, cfg config.Config) error {
	ctx, cancel := context.WithTimeout(parent, cfg.DB.PingTimeout)
	defer cancel()

	return db.PingContext(ctx)
}

func buildLogger(cfg config.Config) logger.Interface {
	if cfg.App.Env == "prod" {
		return logger.Default.LogMode(logger.Warn)
	}
	return logger.Default.LogMode(logger.Info)
}
