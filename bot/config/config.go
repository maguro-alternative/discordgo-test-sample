package config

import (
	"errors"
	"fmt"
	"sync"

	"maguro-alternative/discordgo-test-sample/bot/config/internal"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

var (
	once sync.Once
	cfg  *internal.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	once.Do(MustInit)
}

func MustInit() {
	cfg = &internal.Config{}
	if err := env.Parse(cfg); err != nil {
		xerr := errors.New("failed to parse config: " + err.Error())
		fmt.Printf("panic: %+v", xerr)
		panic(xerr)
	}
}

func DiscordBotToken() string {
	return cfg.DiscordBotToken
}
