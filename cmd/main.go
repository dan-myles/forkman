package main

import (
	"github.com/avvo-na/devil-guard/sentinel"
	"github.com/avvo-na/devil-guard/utils"
)

func main() {
	// Init config and logger
	utils.InitConfig()
	utils.InitLogger()

	// Start the bot
	sentinel.Init()
}
