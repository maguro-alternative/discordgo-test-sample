package main

import (
	"maguro-alternative/discordgo-test-sample/bot"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}

	// ボットの起動
	discord, cleanupCommandHandlers, err := bot.BotOnReady(&client)
	if err != nil {
		panic(err)
	}
	defer cleanupCommandHandlers()
	defer discord.Close()

	// Ctrl+Cを受け取るためのチャンネル
	sc := make(chan os.Signal, 1)
	// Ctrl+Cを受け取る
	signal.Notify(sc, os.Interrupt)
	<-sc //プログラムが終了しないようロック
}
