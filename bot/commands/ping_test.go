package commands

import (
	"testing"

	"maguro-alternative/discordgo-test-sample/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_PingCommand(t *testing.T) {
	session := mock.SessionMock{}
	t.Run("ping成功", func(t *testing.T){
		err := PingCommand().Executor(&session, &discordgo.InteractionCreate{
			Interaction: &discordgo.Interaction{
				Data: discordgo.ApplicationCommandInteractionData{
					Name: "ping",
				},
			},
		})
		assert.NoError(t, err)
	})
}