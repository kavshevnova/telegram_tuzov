package tgClient

import (
	"encoding/json"
	"github.com/kavshevnova/telegrambot_tuzov/lib/e"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string      //хост api сервиса телеграма
	basePath string      //базовый путь(префикс с которого начинаются все запросы)
	client   http.Client // структура  предназначенная для выполнения HTTP-запросов
}

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (client *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) // Добавляет значение к параметру
	q.Add("limit", strconv.Itoa(limit))

	data, err := client.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse
	if err := json.Unmarshal(data, &res); err != nil {
	}
	return res.Result, err
}

func (client *Client) SendMessage(chatId int, message string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", message)
	_, err := client.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}
	return nil
}

func (client *Client) doRequest(method string, q url.Values) ([]byte, error) {
	//request - запрос

	u := url.URL{
		//создаем url запроса для req (запроса NewRequest)
		Scheme: "https",
		Host:   client.host,
		Path:   path.Join(client.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	//создания нового HTTP-запроса. Возвращает *http.Request и ошибку.
	if err != nil {
		return nil, e.WrapIfErr("can't do request", err)
	}
	req.URL.RawQuery = q.Encode() //Кодирование в строку
	//req.URL — это структура url.URL, хранящая разобранный URL запроса.
	//Поле RawQuery содержит параметры после ? (например, ?page=1&sort=asc).
	/*req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
	q := url.Values{}
	q.Add("page", "1")
	q.Add("sort", "desc")
	req.URL.RawQuery = q.Encode() Теперь URL = "https://api.example.com/data?page=1&sort=desc"
	*/
	resp, err := client.client.Do(req)
	//отправляем запрос req запрос и получаем ответ resp
	if err != nil {
		return nil, e.WrapIfErr("can't do request", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body) //Читает все данные из тела HTTP-ответа (resp.Body) и возвращает их в виде байтового среза ([]byte).
	//a) Текст (JSON, HTML, XML): Данные читаются как есть. b) Бинарные данные (фото, видео, PDF): Получаем "сырые" байты.
	if err != nil {
		return nil, e.WrapIfErr("can't do request", err)
	}
	return body, nil
}
