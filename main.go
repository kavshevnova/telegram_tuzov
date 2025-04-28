package main

import (
	"flag"
	tgClient "github.com/kavshevnova/telegrambot_tuzov/clients/telegram"
	event_consumer "github.com/kavshevnova/telegrambot_tuzov/consumer/event-consumer"
	"github.com/kavshevnova/telegrambot_tuzov/events/telegram"
	"github.com/kavshevnova/telegrambot_tuzov/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.NewClient(tgBotHost, musttoken()),
		files.New(storagePath),
	)

	log.Println("Starting telegram bot")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func musttoken() string {
	//приставка маст значит, что вместо ошибки это функция аварийно завершает программу
	token := flag.String("token",
		"",
		"токен тгшки")
	flag.Parse()
	if *token == "" {
		log.Fatal("токена нема")
	}
	return *token
}
