package telegram

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
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
	if len(update.Message.Text) < len("/echo")+1 {
		return
	}
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

func (tgb *tgbot) RaitingHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username == "SlavaYourWarrior" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Слава кыш",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	msg := update.Message.Text[len("/rate")+1:]
	words := strings.Split(msg, " ")
	if len(words) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Неверный запрос",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	if words[0][0] != '@' || len(words[0]) == 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Неверный никнейм",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	username := words[0][1:]

	incStr := words[1]
	inc, err := strconv.Atoi(incStr)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Неверное значение дельты рейтинга",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	val, err := tgb.db.ChangeRaiting(username, inc)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Ошибка: " + err.Error(),
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             "Новый рейтинг пользователя " + username + " " + strconv.Itoa(val),
		ReplyToMessageID: update.Message.ID,
	})
}

func (tgb *tgbot) PingHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.Username != "geogreck" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Антикризисные меры на стадии разработки",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	msg := update.Message.Text[len("/ping")+1:]
	words := strings.Split(msg, " ")
	if len(words) < 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Неверный запрос",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	if words[0][0] != '@' || len(words[0]) == 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Неверный никнейм",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	username := words[0][1:]

	incStr := words[1]
	inc, err := strconv.Atoi(incStr)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "Неверное значение количества пингов",
			ReplyToMessageID: update.Message.ID,
		})
		return
	}

	for i := 0; i < inc; i++ {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           update.Message.Chat.ID,
			Text:             "@" + username,
			ReplyToMessageID: update.Message.ID,
		})
		time.Sleep(time.Second * 1)
	}
}
