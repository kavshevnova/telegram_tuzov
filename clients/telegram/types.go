package telegram

type UpdateResponse struct {
	//response - ответ
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}
type Update struct {
	// update - обновление
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"` //по указателю так как указатель может быть nil то есть сообщение может отсутствовать
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
