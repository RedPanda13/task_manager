package app

import (
	"context"
	"fmt"
	"io"

	"github.com/RedPanda13/task_manager/internal/config"
	"github.com/RedPanda13/task_manager/internal/database"
	"github.com/RedPanda13/task_manager/internal/routes"
)

type container struct {
	server  *config.Server
	closers []io.Closer
}

func (c *container) Close() error {
	var firstErr error

	for _, closer := range c.closers {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func buildContainer(ctx context.Context) (*container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	pg, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}

	healthRoutes := routes.NewHealthRoutes()

	server := config.NewServer(cfg, healthRoutes)

	return &container{
		server: server,
		closers: []io.Closer{
			pg,
		},
	}, nil
}
