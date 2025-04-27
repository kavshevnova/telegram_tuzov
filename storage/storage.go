package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/kavshevnova/telegrambot_tuzov/lib/e"
	"io"
)

type Storage interface {
	//storage - место хранения
	Save(p *Page) error                        //сохранять страницу по ссылке
	PickRandom(userName string) (*Page, error) //возвращать страницу пользователю
	Remove(p *Page) error                      //удалить
	IsExists(p *Page) (bool, error)            //существует ли та или иная страница
}

type Page struct {
	//ссылка, которую мы скинули боту
	URL string
	//имя пользователя, который ее скинул
	Username string
	// Created time.Time чтобы ссылки отправлялись не рандомно
}

var ErrNoSavedPages = errors.New("no saved pages")

func (p Page) Hash() (string, error) {
	//Инициализируем хеш.Хеш (Hash) — это результат работы хеш-функции, которая преобразует произвольные данные (строку, файл, массив байтов) в фиксированную строку символов (обычно шестнадцатеричную) определённой длины.
	//Как работает хеш-функция?
	//Принимает на вход данные любого размера (например, строку "Hello").
	//Преобразует их в последовательность байтов.
	//Вычисляет хеш по специальному алгоритму (SHA-1, SHA-256, MD5 и др.).
	//Возвращает фиксированную строку (например, для SHA-1: "f7ff9e8b7...").
	//Пример для строки "Hello":
	//SHA-1: "f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
	h := sha1.New()
	//io.WriteString записывает строку p.URL в хеш-объект h
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	//возвращаем текстовое представление хеша
	return fmt.Sprintf("%x", h.Sum(nil)), nil
	//данная функция нужна чтобы все файлы имели уникальное имя. Так как файлы с одинаковыми именами нельзя сохранять в одной папке.
}
