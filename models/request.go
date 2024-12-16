package models

// Request - структура для запроса
type Request struct {
	Expression string `json:"expression"`
}

// Response - структура для ответа
type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}
