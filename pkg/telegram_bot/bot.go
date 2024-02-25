package bot

import (
	"log/slog"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramBot struct {
	botAPI  	*tgbotapi.BotAPI
	updates 	tgbotapi.UpdatesChannel
	notify  	chan error
	alert_chat 	string
}

func NewTelegramBot(token string, alert_chat string) *telegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Authorized", "account", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return &telegramBot{
		botAPI:  bot,
		updates: updates,
		notify:  make(chan error, 1),
		alert_chat: alert_chat,
	}
}

func (t *telegramBot) GetUpdates() tgbotapi.UpdatesChannel {
	return t.updates
}

func (t *telegramBot) SendAllert(inputError error) {
	chat_id, err := strconv.Atoi(t.alert_chat)
	if err != nil {
		slog.Error("Pass incorrect alert chat_id", "allert_chat_id", t.alert_chat)
	}
	msg := tgbotapi.NewMessage(int64(chat_id), inputError.Error())
	msg.DisableNotification = true
	_, err = t.botAPI.Send(msg)
	if err != nil {
		slog.Error(
			"The error cannot be transmitted",
			"error", err.Error(),
			"input_error", inputError.Error(),
		)
	}
}

func (t *telegramBot) ReplyToMessage(chatID int64, messageID int, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	msg.DisableNotification = true
	_, err := t.botAPI.Send(msg)
	if err != nil {
		t.SendAllert(err)
	}
	return err
}

func (t *telegramBot) GetFile(fileID string) (string, error) {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	fileLink, err := t.botAPI.GetFile(fileConfig)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	link := fileLink.Link(t.botAPI.Token)
	return link, nil
}

// func (t *telegramBot) Notify() <-chan error {
// 	return t.notify
// }
