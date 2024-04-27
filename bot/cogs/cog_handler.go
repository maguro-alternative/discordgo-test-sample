package cogs

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"

)

type cogHandler struct {
	client *http.Client
}

func newCogHandler(
	client *http.Client,
) *cogHandler {
	return &cogHandler{
		client: client,
	}
}

func RegisterHandlers(
	s *discordgo.Session,
	client *http.Client,
) {
	cogs := newCogHandler(client)
	fmt.Println(s.State.User.Username + "としてログインしました")
	s.AddHandler(cogs.onMessageCreate)
}
