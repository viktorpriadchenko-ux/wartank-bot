// package bot_http_conn -- HTTP-соединение бота
package bot_http_conn

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/bot/bot_net/bot_cookie"
)

// БотХттпСоед -- сетевое соединение бота
type БотХттпСоед struct {
	клиент *http.Client       // Фактический клиент бота
	куки   bot_cookie.БотКуки // Кукисы бота
}

// НовБотХттпСоед -- возвращает новое сетевое соединение бота
func НовБотХттпСоед() *БотХттпСоед {
	клиент := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:      10,
			IdleConnTimeout:   30 * time.Second,
			DisableKeepAlives: false,
		},
		Jar:     nil,
		Timeout: time.Second * 10,
	}
	сам := &БотХттпСоед{
		клиент: клиент,
		куки:   bot_cookie.НовБотКуки(),
	}
	сам.клиент.Jar, _ = cookiejar.New(nil)
	return сам
}

// Клиент -- возвращает сетевого клиента
func (сам *БотХттпСоед) ХттпКлиент() *http.Client {
	return сам.клиент
}

// Куки -- возвращает куки клиента
func (сам *БотХттпСоед) Куки() ИБотКуки {
	return &сам.куки
}
