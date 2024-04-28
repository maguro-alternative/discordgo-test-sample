package commands

import (
	"fmt"

	"maguro-alternative/discordgo-test-sample/testutil/mock"

	"github.com/bwmarrin/discordgo"
)

// スラッシュコマンド内でもデータベースを使用できるようにする
type commandHandler struct{}

func newCogHandler() *commandHandler {
	return &commandHandler{}
}

type command struct {
	Name        string
	Aliases     []string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	AppCommand  *discordgo.ApplicationCommand
	Executor    func(s mock.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error)
}

func (c *command) addApplicationCommand(appCmd *discordgo.ApplicationCommand) {
	c.AppCommand = appCmd
}

type handler struct {
	session  mock.Session
	state    *discordgo.State
	commands map[string]*command
	guild    string
}

// スラッシュコマンドの作成
func newCommandHandler(
	session mock.Session,
	state *discordgo.State,
	guildID string,
) *handler {
	return &handler{
		session:  session,
		state:    state,
		commands: make(map[string]*command),
		guild:    guildID,
	}
}

// スラッシュコマンドの登録
func (h *handler) commandRegister(command *command) error {
	// すでに同じ名前のコマンドが登録されている場合はエラーを返す
	if _, exists := h.commands[command.Name]; exists {
		return fmt.Errorf("command with name `%s` already exists", command.Name)
	}

	fmt.Println(command.Name, "command registered")

	// スラッシュコマンドを登録
	appCmd, err := h.session.ApplicationCommandCreate(
		h.state.User.ID,
		h.guild,
		&discordgo.ApplicationCommand{
			ApplicationID: h.state.User.ID,
			Name:          command.Name,
			Description:   command.Description,
			Options:       command.Options,
		},
	)
	if err != nil {
		return err
	}

	// コマンドに登録したスラッシュコマンドを追加
	command.addApplicationCommand(appCmd)

	// コマンドを登録
	h.commands[command.Name] = command

	// スラッシュコマンドのハンドラを登録
	h.session.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			command.Executor(s, i)
		},
	)

	return nil
}

// スラッシュコマンドの削除
func (h *handler) commandRemove(command *command) error {
	// コマンドが登録されていない場合はエラーを返す
	if command == nil {
		return fmt.Errorf("command does not exist")
	}
	err := h.session.ApplicationCommandDelete(h.state.User.ID, h.guild, command.AppCommand.ID)
	if err != nil {
		return fmt.Errorf("error while deleting application command: %v", err)
	}

	delete(h.commands, command.Name)

	return nil
}

// スラッシュコマンドの取得
func (h *handler) getCommands() []*command {
	var commands []*command

	for _, v := range h.commands {
		commands = append(commands, v)
	}

	return commands
}

func RegisterCommands(discordSession *discordgo.Session) (func(), error) {
	var commandHandlers []*handler
	// 所属しているサーバすべてにスラッシュコマンドを追加する
	// NewCommandHandlerの第二引数を空にすることで、グローバルでの使用を許可する
	commandHandler := newCommandHandler(discordSession, discordSession.State, "")
	// 追加したいコマンドをここに追加
	err := commandHandler.commandRegister(pingCommand())
	if err != nil {
		fmt.Printf("error while registering command: %v\n", err)
		return nil, err
	}
	commandHandlers = append(commandHandlers, commandHandler)
	cleanupCommandHandlers := func() {
		for _, handler := range commandHandlers {
			for _, command := range handler.getCommands() {
				err := handler.commandRemove(command)
				if err != nil {
					fmt.Printf("error while removing command: %v\n", err)
				}
			}
		}
	}
	fmt.Println("commands registered")
	return cleanupCommandHandlers, nil
}
