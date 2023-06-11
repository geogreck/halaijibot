package telegram

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

var (
	prevStatus = false
)

type pingServerResponse struct {
	Online bool `json:"online"`
}

func getServerMessage(ok bool) string {
	if ok {
		return "Minecraft server: ✅Online\n"
	}
	return "Minecraft server: 🟥Offline\n"
}

func (tgb *tgbot) RunServerDaemon(ctx context.Context, addr string, tick time.Duration, chatId string) {
	timer := time.NewTicker(tick)
	tgb.logger.Info("Starting server daemon", zap.Int("seconds", int(tick.Seconds())))
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			tgb.logger.Info("Checking minecraft server status")
			resp, err := http.Get("https://api.mcsrvstat.us/2/" + addr)
			if err != nil {
				tgb.logger.Error("failed to get server status", zap.Error(err))
			}

			var status pingServerResponse
			err = json.NewDecoder(resp.Body).Decode(&status)
			if err != nil {
				tgb.logger.Error("failed to parse response body", zap.Error(err))
			}

			if prevStatus != status.Online {
				tgb.logger.Info("Status changed, sending notification")
				prevStatus = status.Online
				tgb.b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatId,
					Text:   getServerMessage(status.Online),
				})
			} else {

				tgb.logger.Info("Status unchanged")
			}
		}
	}
}
