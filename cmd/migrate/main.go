package main

import (
	"context"
	"go-template/internal/core/config"
	"go-template/internal/db/schema"
	"go-template/internal/log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Error(context.Background(), "cannot load config", "err", err)
	}
	// Context with log
	ctx := log.NewContext(context.Background(), cfg.Log.Level, cfg.Log.Mode, os.Stdout)
	log.Debug(ctx, "database", "url", cfg.Database.URL)

	if err := schema.Migrate(cfg.Database.URL); err != nil {
		log.Error(ctx, "error migrating database", "err", err)
		return
	}

	log.Info(ctx, "migration done!")
}
