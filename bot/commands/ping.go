package commands

import (
	"fmt"

	"maguro-alternative/discordgo-test-sample/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

func pingCommand() *command {
	/*
		pingコマンドの定義

		コマンド名: ping
		説明: Pong!
		オプション: なし
	*/
	exec := newCogHandler()
	return &command{
		Name:        "ping",
		Description: "Pong!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "response",
				Description: "レスポンスのテキストを変更します。",
			},
		},
		Executor: exec.handlePing,
	}
}

func (h *commandHandler) handlePing(s mock.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	/*
		pingコマンドの実行

		コマンドの実行結果を返す
	*/
	if i.Interaction.Data.(discordgo.ApplicationCommandInteractionData).Name != "ping" {
		return nil, nil
	}
	if i.Interaction.GuildID != i.GuildID {
		return nil, nil
	}
	responseMessage := "Pong"
	options := i.ApplicationCommandData().Options
	if len(options) != 0 && options[0].Name == "response" {
		responseMessage = options[0].StringValue()
	}

	commandResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	}
	err := s.InteractionRespond(i.Interaction, commandResponse)
	if err != nil {
		fmt.Printf("error responding to ping command: %v\n", err)
		return nil, err
	}
	return commandResponse, nil
}
