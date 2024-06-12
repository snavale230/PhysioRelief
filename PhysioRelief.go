package main

import (
	"PhysioRelief/helpers"

	"github.com/rs/zerolog/log"
)

func main() {
	server := &helpers.Server{}
	err := server.InitializeRoutes()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize routes")
		return
	}

}
