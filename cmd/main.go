package main

import (
	"github.com/avvo-na/devil-guard/config"
	"github.com/avvo-na/devil-guard/logger"
	"github.com/avvo-na/devil-guard/sentinel"
	"github.com/avvo-na/devil-guard/validator"
)

func init() {
	validator.InitValidator()
	config.InitConfig()
	logger.InitLogger()
	sentinel.InitSentinel()
}

func main() {
	config.WriteEnableModule("utility")

	// s := sentinel.New()
	// s.Start()
	//
	// // Wait here until q is pressed
	// log.Info().Msg("Bot is now running!")
	// log.Info().Msg("Press 'q' to exit")
	//
	// for {
	// 	key, _, err := keyboard.GetSingleKey()
	// 	defer keyboard.Close()
	//
	// 	if err != nil {
	// 		log.Panic().Err(err).Msg("Failed to read key")
	// 	}
	//
	// 	if key == rune('q') {
	// 		log.Info().Msg("Exiting...")
	// 		s.Stop()
	// 		os.Exit(0)
	// 	}
	// }
}
