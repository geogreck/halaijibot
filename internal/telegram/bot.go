package telegram

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/geogreck/halaijibot/internal/storage"
	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

type Bot interface {
	Start(ctx context.Context)
	StartWebhook(ctx context.Context)

	WebhookHandler() http.HandlerFunc
}

type tgbot struct {
	b      *bot.Bot
	logger *zap.Logger
	db     storage.Storage
}

func New(logger *zap.Logger) (Bot, error) {
	token := os.Getenv("TG_TOKEN")
	if token == "" {
		return &tgbot{}, errors.New("no api token specified")
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(DefaultHandler),
		bot.WithDebug(),
	}

	b, err := bot.New(token, opts...)

	if err != nil {
		return &tgbot{}, err
	}

	db, err := storage.New(logger)
	if err != nil {
		logger.Error("failed to create bolt storage", zap.Error(err))
	}

	tgb := &tgbot{
		b:      b,
		logger: logger,
		db:     db,
	}

	hookparams := &bot.SetWebhookParams{
		URL: "bot.grechkogv.ru",
	}
	b.SetWebhook(context.Background(), hookparams)

	inf, _ := b.GetWebhookInfo(context.Background())
	logger.Debug("tg bot webhook:", zap.Any("model", inf))

	tgb.RegisterHandler("/echo", EchoHandler)
	tgb.RegisterHandler("/context", ContextHandler)
	tgb.RegisterHandler("/rate", tgb.RaitingHandler)
	tgb.RegisterHandler("/ping", tgb.PingHandler)

	tgb.RunServerDaemon(context.Background(), "158.160.100.125:25565", time.Minute*11/2, "-1977122895")

	return tgb, nil
}

func (b *tgbot) Start(ctx context.Context) {
	b.b.Start(ctx)
}

func (b *tgbot) StartWebhook(ctx context.Context) {
	b.b.StartWebhook(ctx)
}

func (b *tgbot) WebhookHandler() http.HandlerFunc {
	return b.b.WebhookHandler()
}

func (b *tgbot) RegisterHandler(pattern string, handler bot.HandlerFunc) {
	b.b.RegisterHandler(bot.HandlerTypeMessageText, pattern, bot.MatchTypePrefix, handler)
}
