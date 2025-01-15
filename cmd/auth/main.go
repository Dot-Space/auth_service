package main

import (
	"github.com/Dot-Space/auth_service/config"
)

func main() {
	cfg := config.MustLoad()

	setupLogger(cfg.Env)
}
