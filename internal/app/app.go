package app

import (
	"context"
	"fmt"
)

func Run() error {
	ctx := context.Background()

	c, err := buildContainer(ctx)
	if err != nil {
		return fmt.Errorf("build container: %w", err)
	}
	defer c.Close()

	return c.server.Start()
}
