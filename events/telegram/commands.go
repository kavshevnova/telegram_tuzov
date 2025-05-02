package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/kavshevnova/telegrambot_tuzov/lib/e"
	"github.com/kavshevnova/telegrambot_tuzov/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd    = "/rnd"
	RndAllCmd = "/rnd_all"
	HelpCmd   = "/help"
	StartCmd  = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	//мы будем смотреть на текст сообщения и по его формату и содержанию понимать какая это команда
	text = strings.TrimSpace(text) //удаляем лишние пробелы
	//выводим в консоль кто что пишет в ботика
	log.Printf("got new command '%s' from '%s'", text, username)
	if isAddCmd(text) {
		return p.savePage(text, chatID, username)
	}
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case RndAllCmd:
		return p.sendAll(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendStart(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(pageURL string, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	//подготовим страницу которую собираемся сохранить
	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}
	//смотрим не существует ли уже такая таблица
	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	//отправляем соообщение пользователю если существует
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't do sendRandom", err) }()
	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPage)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendAll(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("Can't do sendAll", err) }()
	pages, err := p.storage.PickAll(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPage)
	}
	var message strings.Builder
	for i, page := range pages {
		message.WriteString(fmt.Sprintf("%d. %s\n", i+1, page.URL))
	}
	if err := p.tg.SendMessage(chatID, message.String()); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendStart(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
