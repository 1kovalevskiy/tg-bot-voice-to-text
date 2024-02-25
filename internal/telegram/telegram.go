package telegram

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/1kovalevskiy/tg-bot-voice-to-text/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	telegramBot interface {
		GetUpdates() tgbotapi.UpdatesChannel
		ReplyToMessage(chatID int64, messageID int, text string) error
		GetFile(fileID string) (string, error)
		SendAllert(inputError error)
	}
	sttAPI interface {
		RecognizeAudio(fileReader io.Reader) (string, error)
	}
)

type loop struct {
	ctx    context.Context
	bot    telegramBot
	stt    sttAPI
	config *config.Config
}

func NewLoop(bot telegramBot, stt sttAPI, config *config.Config) (*loop, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	lp := &loop{
		ctx:    ctx,
		bot:    bot,
		stt:    stt,
		config: config,
	}
	return lp, cancel
}

func (l *loop) Run() error {
	ch := make(chan error, 1)
	go func() {
		ch <- l.run()
	}()
	select {
	case res := <-ch:
		return res
	case <-l.ctx.Done():
		return l.ctx.Err()
	}
}

func (l *loop) run() error {
	slog.Debug("Start loop")
	defer func() error {
		if r := recover(); r != nil {
			return fmt.Errorf("recovered in f: %s", r)
		}
		return nil
	}()

	for update := range l.bot.GetUpdates() {
		slog.Debug("Got update", "update", update)
		if update.Message != nil {
			if !l.validateChat(update.Message) {
				continue
			}
			user_id := update.Message.From.ID
			chat_id := update.Message.Chat.ID
			message_id := update.Message.MessageID
			slog.Debug(
				"Message from allowed chat",
				"chat_username", update.Message.Chat.UserName,
				"username", update.Message.From.UserName,
			)

			if update.Message.Voice == nil {
				continue
			}
			voice_id := update.Message.Voice.FileID
			slog.Debug("Got voice message", "voice", update.Message.Voice)

			slog.Info(
				"Got message with Voice",
				"user", user_id,
				"chat", chat_id,
				"voice_file_id", voice_id,
			)
			link, err := l.bot.GetFile(voice_id)
			if err != nil {
				continue
			}
			go func() {
				link := link
				chat_id := chat_id
				message_id := message_id
				resp, err := http.Get(link)
				if err != nil {
					slog.Error(err.Error())
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					slog.Error("bad status", "status", resp.Status)
					return
				}
				text, err := l.stt.RecognizeAudio(resp.Body)
				if err != nil {
					l.bot.SendAllert(err)
				}
				l.bot.ReplyToMessage(chat_id, message_id, text)
			}()

		}
	}
	return nil
}

func (l *loop) validateChat(message *tgbotapi.Message) bool {
	chat := message.Chat
	if chat == nil {
		return false
	}
	allowChats := l.config.AllowChats
	for _, allowChat := range allowChats {
		if strconv.Itoa(int(chat.ID)) == allowChat {
			return true
		}
	}
	allowUsers := l.config.AllowUsers
	for _, allowUser := range allowUsers {
		if strconv.Itoa(int(chat.ID)) == allowUser {
			return true
		}
	}

	slog.Info(
		"Message from disallowed chat",
		"id", message.Chat.ID,
		"username", message.Chat.Title,
	)
	return false
}
