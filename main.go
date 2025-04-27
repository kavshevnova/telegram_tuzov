package main

import (
	"flag"
	"log"
)

const tgBotHost = "api.telegram.org"

func main() {
	t := musttoken()

	tgClient = telegram.New(musttoken())

	//fetcher = fetcher.New() отправлять запрос

	//processor = processor.New() получать сообщения

	//consumer.Start(fetcher, processor)
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
}
