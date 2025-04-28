package event_consumer

import (
	"github.com/kavshevnova/telegrambot_tuzov/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int //размер пачки(сколько событий мы будем обрабатывать за раз)
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	//здесь будет вечный цикл, который постоянно будет ждать новые события и обрабатывать их
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize) //лучше встроить какой-то механизм летрая
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}
		//если событий, 0 то так же пропускаем итерацию, но ждем 1 секунду
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}
		if err := c.handleEvent(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

/*
1. потеря событий: ретраи, возвращение в хранилище, фоллбэк, подтверждение
2.обработка всей пачки: останавливаться после первой ошибки, счетчик ошибок
3.параллельная обработка
*/
func (c Consumer) handleEvent(events []events.Event) error {
	//перебираем события
	for _, event := range events {
		//выводим текст события
		log.Printf("got new event: %s", event.Text)
		//если пошло что-то не так выводим и пропускаем обработку
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
