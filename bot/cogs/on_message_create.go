package cogs

import (
	"context"
	"strings"

	"log/slog"
	"net/http"

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
) (*discordgo.Message, error) {
	// Botの発言は無視
	if vs.Author.Bot {
		return nil, nil
	}

	// メッセージが "ping" の場合は、"pong" と返信
	if vs.Message.Content == "ping" {
		_, err := s.ChannelMessageSend(vs.ChannelID, "pong", discordgo.WithClient(client))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	// !hello が含まれている場合は、何もしない
	if strings.Contains(vs.Message.Content, "!hello") {
		return nil, nil
	}

	// メッセージが自分の発言でない場合は、メッセージを返信
	message, err := s.ChannelMessageSend(vs.ChannelID, "Hello, World!", discordgo.WithClient(client))
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "Message sent", "Message:", message.ID, "Content:", message.Content)

	return message, nil
}
