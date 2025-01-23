package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kpechenenko/botcleaner/internal/bot"
	"github.com/kpechenenko/botcleaner/internal/cache/lru"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))
	cfg, err := loadFromEnv()
	if err != nil {
		logger.Error("load config", err)
		return
	}
	c, err := lru.New[string, int](cfg.cacheCapacity)
	if err != nil {
		logger.Error("create lru cache", err)
		return
	}
	p := bot.CreateBotParams{
		AlertMessageTemplate: cfg.alertMessageTemplate,
		TrackedChannels:      cfg.trackedChannels,
		TgBotToken:           cfg.tgBotToken,
	}
	tgBot, err := bot.New(p, c, logger)
	if err != nil {
		logger.Error("create tg bot:", err)
		return
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		<-signals
		logger.Info("bot shutting down")
		tgBot.StopPolling()
		close(done)
		logger.Info("bot shut down")
	}()
	go func() {
		logger.Info("starting bot")
		tgBot.StartPolling()
	}()
	<-done
}
