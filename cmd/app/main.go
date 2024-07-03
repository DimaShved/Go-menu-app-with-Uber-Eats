package main

import (
	"os"
	"uber-go-menu-copy/internal/config"
	"uber-go-menu-copy/internal/pkg/db"
)

func main() {
	cfg := config.LoadConfig()
	err := db.Connect(&cfg.Database)
	if err != nil {
		os.Exit(1)
	}
}
