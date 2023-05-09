package telegram

import (
	"context"
	"errors"
	"net/http"
	"os"

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
}

func New(logger *zap.Logger) (Bot, error) {
	token := os.Getenv("TG_TOKEN")
	if token == "" {
		return &tgbot{}, errors.New("no api token specified")
	}

	opts := []bot.Option{}

	b, err := bot.New(token, opts...)

	if err != nil {
		return &tgbot{}, err
	}

	tgb := &tgbot{
		b:      b,
		logger: logger,
	}

	tgb.RegisterHandler("/echo", EchoHandler)

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
	b.b.RegisterHandler(bot.HandlerTypeMessageText, pattern, bot.MatchTypeExact, handler)
}