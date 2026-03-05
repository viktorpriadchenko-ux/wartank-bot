package types

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
)

// ИБотоФерма -- словарь серверных ботов
type ИБотоФерма interface {
	// Get -- возвращает бота по его имени
	Get(botNumber АБотНомер) ИБот
	// BotStart -- запускает бота в работу
	BotStart(botNumber АБотНомер) Result[bool]
	// ListBot -- возвращает список ботов
	ListBot() []ИБот
	// НовБот -- добавляет нового бота на бото-ферму
	НовБот(login, password string, еслиАвто bool) Result[bool]
}
