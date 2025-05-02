package main

import (
	"context"
	"flag"
	tgClient "github.com/kavshevnova/telegrambot_tuzov/clients/telegram"
	event_consumer "github.com/kavshevnova/telegrambot_tuzov/consumer/event-consumer"
	"github.com/kavshevnova/telegrambot_tuzov/events/telegram"
	"github.com/kavshevnova/telegrambot_tuzov/storage/sqlite"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "storage/sqlite/storage.db"
	batchSize         = 100
)

func main() {

	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect sqlite storage: ", err)
	}
	//используя туду мы говорим что мы еще не определились с тем контекстом который мы будем использовать и потом мы можем здесь же его заменить на другой
	//если мы используем контекст бэкграунд мы четко говорим что здесь нужен контекст который никак нас не ограничивает, дальше от него будет унаследован другой контекст типо дедлайн или таймаут
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init sqlite storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.NewClient(tgBotHost, musttoken()),
		s,
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
