package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/1kovalevskiy/tg-bot-voice-to-text/config"
	"github.com/1kovalevskiy/tg-bot-voice-to-text/internal/telegram"
	yandexcloudstt "github.com/1kovalevskiy/tg-bot-voice-to-text/internal/yandex_cloud_stt"
	"github.com/1kovalevskiy/tg-bot-voice-to-text/pkg/logger"
	bot "github.com/1kovalevskiy/tg-bot-voice-to-text/pkg/telegram_bot"
	"github.com/1kovalevskiy/tg-bot-voice-to-text/pkg/yandex_cloud"
)

func Run(cfg *config.Config) {
	const op = "internal - app - Run"
	logger.New(cfg.Env.Environment)
	yandex_client := yandex_cloud.NewgRPCClient(cfg.YandexOauth)
	defer yandex_client.Shutdown()

	bot := bot.NewTelegramBot(cfg.TelegramToken, cfg.AlertChat)

	sst_client := yandexcloudstt.NewSSTClient(yandex_client)

	loop, cancel := telegram.NewLoop(bot, sst_client, cfg)

	ch := make(chan error, 1)
	go func() {
		ch <- loop.Run()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		slog.Info(op + " - signal: " + s.String())
	case err := <-ch:
		slog.Error(op+" - telegram.Run:", "error", err.Error())
	}

	cancel()

}
