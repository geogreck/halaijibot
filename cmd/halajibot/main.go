package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/geogreck/halaijibot/internal/telegram"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bot, err := telegram.New(logger)
	if err != nil {
		logger.Fatal("Failed to create tg bot", zap.Error(err))
	}

	logger.Info("Created bot successfully")

	go bot.StartWebhook(context.Background())

	logger.Info("Started webhook")

	logger.Info("Starting server")
	http.ListenAndServe("0.0.0.0:8080", bot.WebhookHandler())
}
