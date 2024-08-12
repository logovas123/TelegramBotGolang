package main

import (
	"context"
	"flag"
	"log"

	tgClient "TelegramBotGolang/clients/telegram"
	"TelegramBotGolang/consumer/eventConsumer"
	"TelegramBotGolang/events/telegram"
	"TelegramBotGolang/storage/sqlite"
)

const (
	tgBotHost = "api.telegram.org"
	// всё храним в папке data
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// s := files.New(storagePath)
	// создаём storage
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err = s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), s)

	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
