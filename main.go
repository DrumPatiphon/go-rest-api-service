package main

import (
	"os"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/severs"
	"github.com/DrumPatiphon/go-rest-api-service/pkg/database"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DbConnect(cfg.Db())
	defer db.Close()

	severs.NewSever(cfg, db).Start()
}
