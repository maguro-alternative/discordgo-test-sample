package commands

import (
	"errors"
	"testing"

	"maguro-alternative/discordgo-test-sample/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_pingCommand(t *testing.T) {
	session := mock.SessionMock{
		InteractionRespondFunc: func(i *discordgo.Interaction, r *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
			return nil
		},
	}
	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Data: discordgo.ApplicationCommandInteractionData{
				Name: "ping",
			},
		},
	}
	interaction.Type = discordgo.InteractionApplicationCommand
	interaction.GuildID = "1234567890"
	t.Run("ping成功", func(t *testing.T) {
		interaction.Interaction.GuildID = "1234567890"
		commandResponse, err := pingCommand().Executor(&session, interaction)
		assert.NoError(t, err)
		assert.NotNil(t, commandResponse)
		assert.Equal(t, "Pong", commandResponse.Data.Content)
	})

	t.Run("ping以外のスラッシュコマンドの場合終了", func(t *testing.T) {
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "pong",
				},
			},
		}
		interaction.Type = discordgo.InteractionApplicationCommand
		interaction.Interaction.GuildID = "1234567890"
		session := mock.SessionMock{}
		commandResponse, err := pingCommand().Executor(&session, interaction)
		assert.NoError(t, err)
		assert.Nil(t, commandResponse)
	})

	t.Run("ping失敗", func(t *testing.T) {
		interaction.Interaction.GuildID = "1234567890"
		session := mock.SessionMock{
			InteractionRespondFunc: func(i *discordgo.Interaction, r *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
				return errors.New("error")
			},
		}
		commandResponse, err := pingCommand().Executor(&session, interaction)
		assert.Error(t, err)
		assert.Nil(t, commandResponse)
	})

	t.Run("pingのオプション取得", func(t *testing.T) {
		interaction := &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "ping",
					Options: []*discordgo.ApplicationCommandInteractionDataOption{
						{
							Type:  discordgo.ApplicationCommandOptionString,
							Name:  "response",
							Value: "Pong!!!!",
						},
					},
				},
			},
		}
		interaction.Type = discordgo.InteractionApplicationCommand
		interaction.Interaction.GuildID = "1234567890"
		commandResponse, err := pingCommand().Executor(&session, interaction)
		assert.NoError(t, err)
		assert.NotNil(t, commandResponse)
		assert.Equal(t, "Pong!!!!", commandResponse.Data.Content)
	})
}
