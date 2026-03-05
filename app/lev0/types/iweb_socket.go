package types

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ИВебСокет -- интерфейс к постоянному веб-сокету сервера
type ИВебСокет interface {
	// Записать -- записывает топик на сервер
	Записать(topic string, dictReq map[string]string) error
	// Читать -- читает топик с сервера
	Читать(topic string) (mapResp map[string]string, err error)
	// Вызвать -- вызывает удалённую процедуру
	Вызвать(topic string, dictReq map[string]string) Result[map[string]string]
	// ЕслиПодключ -- возвращает признак подключенности к серверу
	ЕслиПодключ() bool
}
