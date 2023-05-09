package telegram

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Контекст не ясен",
	})
}

func EchoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text[len("/echo")+1:],
	})
}

func ContextHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	i := rand.New(rand.NewSource(time.Now().Unix())).Intn(2)
	if i == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Контекст ясен",
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Контекст не ясен",
		})
	}
}

func RaitingHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username == "SlavaYourWarrior" || update.Message.From.Username == "geogreck" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Слава кыш",
			ReplyToMessageID: update.Message.ID,
		})
	}
	// msg := update.Message.Text[len("/rate"):]
}
