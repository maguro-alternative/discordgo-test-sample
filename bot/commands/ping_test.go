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
	interaction.GuildID = "1234567890"
	t.Run("ping成功", func(t *testing.T) {
		interaction.Interaction.GuildID = "1234567890"
		err := pingCommand().Executor(&session, interaction)
		assert.NoError(t, err)
	})

	t.Run("ping失敗", func(t *testing.T) {
		interaction.Interaction.GuildID = "1234567890"
		session := mock.SessionMock{
			InteractionRespondFunc: func(i *discordgo.Interaction, r *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
				return errors.New("error")
			},
		}
		err := pingCommand().Executor(&session, interaction)
		assert.Error(t, err)
	})
}
