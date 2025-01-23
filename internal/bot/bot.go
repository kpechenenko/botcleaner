package bot

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kpechenenko/botcleaner/internal/cache"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	tb                   *telego.Bot
	cache                cache.Cache[string, int]
	trackedChannels      map[string]bool
	alertMessageTemplate string
	stop                 chan struct{}
	logger               *slog.Logger
}

func (b *Bot) StartPolling() {
	updates, _ := b.tb.UpdatesViaLongPolling(nil)
	for {
		select {
		case <-b.stop:
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			if update.Message == nil {
				continue
			}
			switch origin := update.Message.ForwardOrigin.(type) {
			case *telego.MessageOriginChannel:
				b.handleForwardedMeme(update, origin)
			}
		}
	}
}

func (b *Bot) handleForwardedMeme(update telego.Update, origin *telego.MessageOriginChannel) {
	if !b.channelTracked(origin) {
		return
	}
	memeMessageId, ok := b.loadMemeMessageIdFromCache(origin, update)
	if !ok {
		b.addMemeMessageIdToCache(origin, update)
		return
	}
	_, _ = b.sendAlertAboutRepeatedMeme(update, memeMessageId)
	_ = b.deleteRepeatedMemeMessage(update)
}

func (b *Bot) sendAlertAboutRepeatedMeme(update telego.Update, memeMessageId int) (*telego.Message, error) {
	text := fmt.Sprintf(b.alertMessageTemplate, update.Message.From.Username)
	msg, err := b.tb.SendMessage(&telego.SendMessageParams{
		ChatID: tu.ID(update.Message.Chat.ID),
		Text:   text,
		ReplyParameters: &telego.ReplyParameters{
			MessageID:                memeMessageId,
			ChatID:                   tu.ID(update.Message.Chat.ID),
			AllowSendingWithoutReply: true,
		},
	})
	if err != nil {
		b.logger.Error(fmt.Sprintf("send alert: %v", err))
	}
	return msg, err
}

func (b *Bot) deleteRepeatedMemeMessage(update telego.Update) error {
	err := b.tb.DeleteMessage(&telego.DeleteMessageParams{
		MessageID: update.Message.MessageID,
		ChatID:    tu.ID(update.Message.Chat.ID),
	})
	if err != nil {
		b.logger.Error(fmt.Sprintf("delete message: %v", err))
	}
	return err
}

func (b *Bot) channelTracked(origin *telego.MessageOriginChannel) bool {
	return b.trackedChannels[origin.Chat.Username]
}

func (b *Bot) generateCacheKeyToStoreMeme(origin *telego.MessageOriginChannel, update telego.Update) string {
	// key structure: chatId of received chat $ chatId of src tg channel $ messageId in src tg channel
	return fmt.Sprintf("%d$%d$%d", update.Message.Chat.ID, origin.Chat.ID, origin.MessageID)
}

func (b *Bot) loadMemeMessageIdFromCache(origin *telego.MessageOriginChannel, update telego.Update) (int, bool) {
	messageId, ok := b.cache.Get(b.generateCacheKeyToStoreMeme(origin, update))
	return messageId, ok
}

func (b *Bot) addMemeMessageIdToCache(origin *telego.MessageOriginChannel, update telego.Update) {
	b.cache.Set(b.generateCacheKeyToStoreMeme(origin, update), update.Message.MessageID)
}

func (b *Bot) StopPolling() {
	b.tb.StopLongPolling()
	close(b.stop)
}

type CreateBotParams struct {
	AlertMessageTemplate string
	TrackedChannels      []string
	TgBotToken           string
}

func (p *CreateBotParams) check() error {
	if p.AlertMessageTemplate == "" {
		return errors.New("alert text template can't be empty")
	}
	if len(p.TrackedChannels) == 0 {
		return errors.New("tracked channels can't be empty")
	}
	if p.TgBotToken == "" {
		return errors.New("bot token can't be empty")
	}
	return nil
}

func extractBaseTgChannelName(name string) string {
	return strings.TrimPrefix(
		strings.TrimSpace(name),
		"@",
	)
}

func New(params CreateBotParams, cache cache.Cache[string, int], log *slog.Logger) (*Bot, error) {
	if err := params.check(); err != nil {
		log.Error(fmt.Sprintf("check params: %v", err))
		return nil, err
	}
	tb, err := telego.NewBot(params.TgBotToken, telego.WithLogger(&logger{log: log}))
	if err != nil {
		log.Error(fmt.Sprintf("create telego bot: %v", err))
		return nil, err
	}
	return &Bot{
		tb:                   tb,
		cache:                cache,
		alertMessageTemplate: params.AlertMessageTemplate,
		trackedChannels:      convertListToSet(params.TrackedChannels, extractBaseTgChannelName),
		stop:                 make(chan struct{}, 1),
	}, nil
}
