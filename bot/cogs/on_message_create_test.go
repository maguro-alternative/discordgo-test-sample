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

func TestLineRequest_PushMessageNotify(t *testing.T) {
	ctx := context.Background()
	// スタブHTTPクライアントを作成
	stubClient := mock.NewStubHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}
	})
	t.Run("正常系", func(t *testing.T) {
		discordState := discordgo.NewState()
		discordState.User = &discordgo.User{
			ID:       "111",
			Username: "test",
		}
		content, err := onMessageCreateFunc(
			ctx,
			stubClient,
			&mock.SessionMock{
				ChannelFunc: func(channelID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error) {
					return &discordgo.Channel{
						GuildID: "123",
					}, nil
				},
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{}, nil
				},
				GuildFunc: func(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error) {
					return &discordgo.Guild{
						ID: "123",
					}, nil
				},
			},
			discordState,
			&discordgo.MessageCreate{
				Message: &discordgo.Message{
					Content: "test",
					Author: &discordgo.User{
						Bot: false,
					},
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, "t", content)
	})
}
