package main

import (
	"PhysioRelief/helpers"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ViperInitialization initializes Viper by adding the config path, setting the config name, and reading the config from the file. If there is an error reading the config, it logs it
func ViperInitialization() error {
	viper.AddConfigPath("/app")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("viper config not found")
		return err
	}
	return nil
}

func main() {
	server := &helpers.Server{}
	err := server.InitializeRoutes()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize routes")
		return
	}

}
