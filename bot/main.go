package bot

import (
	"net/http"

	"maguro-alternative/discordgo-test-sample/bot/config"
	"maguro-alternative/discordgo-test-sample/bot/commands"
	"maguro-alternative/discordgo-test-sample/bot/cogs"

	"github.com/cockroachdb/errors"
	"github.com/bwmarrin/discordgo"
)

func BotOnReady(client *http.Client) (*discordgo.Session, func(), error) {
	/*
		ボットの起動

		args:
		indexDB: db.Driver
		データベースのドライバー

		return:
		*discordgo.Session
		エラーがなければ、セッションを返します。
		エラーがあれば、エラーを返します。
	*/
	// セッションを作成
	discordToken := "Bot " + config.DiscordBotToken()
	discordSession, err := discordgo.New(discordToken)
	if err != nil {
		return nil, func(){}, errors.WithStack(err)
	}
	discordSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	discordSession.Token = discordToken
	err = discordSession.Open()
	if err != nil {
		return nil, func(){}, errors.WithStack(err)
	}
	cogs.RegisterHandlers(discordSession, client)
	cleanupCommandHandlers, err := commands.RegisterCommands(discordSession)
	return discordSession, cleanupCommandHandlers, err
}
