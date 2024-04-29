# discordgoのe2eテストサンプル
discordgoのe2eテストのサンプルです。

# 準備
```.env```ファイルを作成し、以下の内容を記述してください。
```
DISCORD_BOT_TOKEN=DiscordのBotのトークン
```

Botの起動は以下のコマンドで行います。
```
go run core/main.go
```

# 共通事項
## 引数と戻り値
discordgoではAddHandler関数を使ってイベント処理を登録します。  
しかし、引数と戻り値の指定がある都合上、戻り値に```error```を指定できず、正しい挙動をしたかどうかの確認ができません。
```go:ok.go
// func(セッション, イベント)の形式でイベント処理を登録
func onMessageCreate(s *discordgo.Session, vs *discordgo.MessageCreate) {}
```

```go:ng.go
// func(セッション, イベント)の形式でないため登録できない
func onMessageCreate(
	ctx context.Context,
	client *http.Client,
	s *discordgo.Session,
	state *discordgo.State,
	vs *discordgo.MessageCreate,
) (*discordgo.Message, error) {}
```

というわけで、以下のように登録する関数内で、処理を行う関数を呼びだします。
```go
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
	s *discordgo.Session,
	state *discordgo.State,
	vs *discordgo.MessageCreate,
) (*discordgo.Message, error) {
    // 処理
}
```

## モック
discordgoの送受信の処理をモック化するために、以下のようなインターフェースを作成します。
```go:testutil/mock/session.go
type SessionMock struct {
    ChannelMessageSendFunc func(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

func (s *SessionMock) ChannelMessageSend(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
    return s.ChannelMessageSendFunc(channelID, content, options...)
}

type Session interface {
    ChannelMessageSend(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
}

var (
	_ Session = (*discordgo.Session)(nil)
	_ Session = (*SessionMock)(nil)
)
```

モックを使用するため、関数の引数を```*discordgo.Session```から```mock.Session```に変更します。

```diff
func onMessageCreateFunc(
	ctx context.Context,
	client *http.Client,
-	s *discordgo.Session,
+	s mock.Session,
	state *discordgo.State,
	vs *discordgo.MessageCreate,
) (*discordgo.Message, error) {
    // 処理
}
```



# cogのテスト
受け取ったイベントに対する処理をテストします。  
サンプルではメッセージが送信された際の処理をテストしています。  

テストする```onMessageCreateFunc```は以下のような仕様です。
- Botからのメッセージが送信された場合、何も返さない。
- ```ping```が送信された場合、```pong```を返す。
- ```!hello```が送信された場合、何も返さない。
- 上記以外の場合、```Hello, World!```を返す。

```go:bot/cogs/on_message_create.go
func onMessageCreateFunc(
	ctx context.Context,
	client *http.Client,
	s mock.Session,
	state *discordgo.State,
	vs *discordgo.MessageCreate,
) (*discordgo.Message, error) {}
```

テストは以下のような形で行います。  
```Hello, World!```が返ってくるかどうかを確認しています。
```go:bot/cogs/on_message_create_test.go
func TestOnMessageCreateFunc(t *testing.T) {
    // テスト対象の関数を呼び出し
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
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						ID:      channelID,
						Content: content,
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
}
```

```&mock.SessionMock```を引数に渡すことで、モックを使用してテストを行います。  
実際に送信はされず、送信されたとしてレスポンスを返すようにしています。  

その他、以下のようなテストも行っています。
- メッセージが送信されたが、Botからのメッセージの場合
- 特定の文字のメッセージが送信された場合
- メッセージの送信に失敗した場合

# スラッシュコマンドのテスト
スラッシュコマンドは登録のテストも行います。  
以下の構造体に、コマンドの情報を保持できるかテストします。
```go
type handler struct {
	session  mock.Session
	state    *discordgo.State
	commands map[string]*command
	guild    string
}
```

```commandRegister```関数でコマンドを登録し、```handler.commands```に格納します。  
```commandRemove```関数でコマンドを削除し、```handler.commands```内のデータも削除します。  
```getCommand```関数で```handler.commands```からコマンドを取得します。  

各スラッシュコマンドは以下のように定義します。
```go
func pingCommand() *command {
	/*
		pingコマンドの定義

		コマンド名: ping
		説明: Pong!
		オプション: なし
	*/
	exec := newCogHandler()
	return &command{
		Name:        "ping",
		Description: "Pong!",
		Options:     []*discordgo.ApplicationCommandOption{
			{
                Required:    false,
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "response",
				Description: "レスポンスのテキストを変更します。",
			},
		},
		Executor:    exec.handlePing,
	}
}

func (h *commandHandler) handlePing(s mock.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {}
```

Optionはスラッシュコマンドのオプション(discord.go公式のサンプルコード参照)、Executorは実行する関数を指定します。  

スラッシュコマンドの関数は以下のような引数と戻り値の指定があります。
```go
func(s mock.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error)
```

コマンドに対する返信の内容とエラーを返します。  

テストは以下のように行います。
```go
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
```

```cog```と同様に```&mock.SessionMock```を引数に渡すことで、モックを使用してテストを行います。  
その他、以下のようなテストも行っています。
- 別のコマンドが実行された場合
- コマンドの返信に失敗した場合
- オプションが指定された場合

必要に応じてテストを追加した下さい。
