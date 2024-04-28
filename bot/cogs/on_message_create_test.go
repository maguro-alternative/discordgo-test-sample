package cogs

import (
	"context"
	"io"
	"maguro-alternative/discordgo-test-sample/testutil/mock"
	"net/http"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestOnMessageCreate(t *testing.T) {
	ctx := context.Background()
	// スタブHTTPクライアントを作成
	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})
	t.Run("正常系(Hello World!を返す)", func(t *testing.T) {
		discordState := discordgo.NewState()
		discordState.User = &discordgo.User{
			ID:       "111",
			Username: "test",
		}
		message, err := onMessageCreateFunc(
			ctx,
			stubClient,
			&mock.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "123",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						ID:      channelID,
						Content: content,
					}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
					return &discordgo.Guild{
						ID: guildID,
					}, nil
				},
			},
			discordState,
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Author: &discordgo.User{
						Bot: false,
					},
					Content: "test",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World!", message.Content)
	})

	t.Run("正常系(pingを受け取った場合pongを返す)", func(t *testing.T) {
		discordState := discordgo.NewState()
		discordState.User = &discordgo.User{
			ID:       "111",
			Username: "test",
		}
		message, err := onMessageCreateFunc(
			ctx,
			stubClient,
			&mock.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "123",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						ID:      channelID,
						Content: content,
					}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error) {
					return &discordgo.Guild{
						ID: guildID,
					}, nil
				},
			},
			discordState,
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Author: &discordgo.User{
						Bot: false,
					},
					Content: "ping",
				},
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "pong", message.Content)
	})
}

