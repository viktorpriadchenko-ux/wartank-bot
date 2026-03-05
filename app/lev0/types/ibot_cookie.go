package types

import (
	"net/http"
)

// ИБотКуки -- куки серверного бота
type ИБотКуки interface {
	// Уст -- устанавливает куки бота
	Уст(куки []*http.Cookie)
	// Получ -- возвращает куки бота
	Получ() []*http.Cookie
}
