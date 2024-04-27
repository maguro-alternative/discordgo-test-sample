package cogs

import (
	"context"

	"log/slog"
	"net/http"

	//"maguro-alternative/discordgo-test-sample/bot/config"
	"maguro-alternative/discordgo-test-sample/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

func (h *cogHandler) onMessageCreate(s *discordgo.Session, vs *discordgo.MessageCreate) {
	ctx := context.Background()
	_, err := onMessageCreateFunc(ctx, h.client, s, s.State, vs)
	if err != nil {
		slog.ErrorContext(ctx, "OnMessageCreate Error", "Error:", err.Error())
	}
}

func onMessageCreateFunc(
	ctx context.Context,
	client *http.Client,
	s mock.Session,
	state *discordgo.State,
	vs *discordgo.MessageCreate,
) (string, error) {
	// メッセージが自分の発言かどうかを判定
	if vs.Author.ID == state.User.ID {
		return "", nil
	}

	// メッセージが自分の発言でない場合は、メッセージを返信
	message, err := s.ChannelMessageSend(vs.ChannelID, "Hello, World!", discordgo.WithClient(client))
	if err != nil {
		return "", err
	}
	slog.InfoContext(ctx, "Message sent", "Message:", message.ID, "Content:", message.Content)

	return message.Content, nil
}
