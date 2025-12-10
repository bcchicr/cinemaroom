package main

import (
	"log"

	"github.com/bcchicr/cinemaroom/services/authn/app"
	"github.com/bcchicr/cinemaroom/services/authn/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	app.Run(cfg)
}
