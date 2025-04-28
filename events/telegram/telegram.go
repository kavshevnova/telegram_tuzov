package telegram

import (
	"errors"
	"github.com/kavshevnova/telegrambot_tuzov/clients/telegram"
	"github.com/kavshevnova/telegrambot_tuzov/events"
	"github.com/kavshevnova/telegrambot_tuzov/lib/e"
	"github.com/kavshevnova/telegrambot_tuzov/storage"
)

type Processor struct {
	tg      *tgClient.Client //тг
	offset  int              //смещение
	storage storage.Storage  //место хранения
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknownEvent = errors.New("unknown event")
var ErrUnknownMetaType = errors.New("unknown meta type")

func New(client *tgClient.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	//мы получаем апдейты
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.WrapIfErr("can't fetch updates", err)
	}
	//если список апдейтов оказался пустым то мы заканчиваем работу функции и говорим что ничего не нашли
	if len(updates) == 0 {
		return nil, nil
	}
	//объявляем переменную для результата и заранее обозначаем память для нее так нам известно сколько будет значений
	res := make([]events.Event, 0, len(updates))
	//перебираем все апдейты и преобразуем их в тип event
	for _, u := range updates {
		res = append(res, event(u))
	}
	//обновляем параметр офсет чтобы в следующий раз получить следующую пачку изменений
	p.offset = updates[len(updates)-1].ID + 1 //при следующем запросе мы получим только те апдейты у которых айди больше чем у последнего из уже полученных
	//возвращаем результат
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	//этот метод будет выполнять различные действия в зависимости от типа ивента
	switch event.Type {
	case events.Message:
		return p.processMessage(event) //функция логики работы с сообщением
	default:
		return e.Wrap("can't process message", ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd tgClient.Update) events.Event {
	//В чем разница между ивентами и апдейтами: апдейты это параметр телеграма и они относятся только к нему, а ивент это более общая сущность, в нее мы можем преобразовывать все что получаем от других мессенжеров, в каком бы формате они не предоставляли нам информацию.
	updType := fetchType(upd)
	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = &Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchText(upd tgClient.Update) string {
	if upd.Message != nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd tgClient.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
