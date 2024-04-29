package commands

import (
	"testing"

	"maguro-alternative/discordgo-test-sample/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_PingCommandRegister(t *testing.T) {
	session := mock.SessionMock{
		AddHandlerFunc: func(interface{}) func() {
			return func() {}
		},
		ApplicationCommandCreateFunc: func(guildID string, appID string, appCmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error) {
			return appCmd, nil
		},
	}
	t.Run("ping登録成功", func(t *testing.T) {
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(pingCommand())
		assert.NoError(t, err)
		c := h.getCommands()
		assert.NotNil(t, c)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "ping", c[0].Name)
	})

	t.Run("ping登録失敗(2重登録)", func(t *testing.T) {
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(pingCommand())
		assert.NoError(t, err)
		c := h.getCommands()
		assert.NotNil(t, c)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "ping", c[0].Name)

		err = h.commandRegister(pingCommand())
		assert.Error(t, err)
	})

	t.Run("ping登録失敗(エラー発生)", func(t *testing.T) {
		session := mock.SessionMock{
			AddHandlerFunc: func(interface{}) func() {
				return func() {}
			},
			ApplicationCommandCreateFunc: func(guildID string, appID string, appCmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error) {
				return nil, assert.AnError
			},
		}
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(pingCommand())
		assert.Error(t, err)

		c := h.getCommands()
		assert.Len(t, c, 0)
	})
}

func Test_PingCommandRemove(t *testing.T) {
	session := mock.SessionMock{
		AddHandlerFunc: func(interface{}) func() {
			return func() {}
		},
		ApplicationCommandCreateFunc: func(guildID string, appID string, appCmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (ccmd *discordgo.ApplicationCommand, err error) {
			return appCmd, nil
		},
	}
	t.Run("ping削除成功", func(t *testing.T) {
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(pingCommand())
		assert.NoError(t, err)
		c := h.getCommands()
		assert.NotNil(t, c)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "ping", c[0].Name)

		err = h.commandRemove(c[0])
		assert.NoError(t, err)

		c = h.getCommands()
		assert.Len(t, c, 0)
	})

	t.Run("ping削除失敗(未登録)", func(t *testing.T) {
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		session.ApplicationCommandDeleteFunc = func(appID string, guildID string, cmdID string, options ...discordgo.RequestOption) error {
			return assert.AnError
		}
		h := newCommandHandler(&session, state, "")
		c := h.getCommands()
		assert.Len(t, c, 0)

		err := h.commandRemove(&command{
			AppCommand: &discordgo.ApplicationCommand{
				ID: "1234567890",
			},
			Name: "ping",
		})
		assert.Error(t, err)
	})

	t.Run("ping削除失敗(エラー発生)", func(t *testing.T) {
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRegister(pingCommand())
		assert.NoError(t, err)
		c := h.getCommands()
		assert.NotNil(t, c)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "ping", c[0].Name)

		session.ApplicationCommandDeleteFunc = func(appID string, guildID string, cmdID string, options ...discordgo.RequestOption) error {
			return assert.AnError
		}
		err = h.commandRemove(c[0])
		assert.Error(t, err)
	})

	t.Run("ping削除失敗(削除対象がnil)", func(t *testing.T) {
		state := discordgo.NewState()
		state.User = &discordgo.User{
			ID: "1234567890",
		}
		h := newCommandHandler(&session, state, "")
		err := h.commandRemove(nil)
		assert.Error(t, err)
	})
}
