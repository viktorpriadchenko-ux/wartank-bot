package types

import (
	. "wartank/app/lev0/alias"
)

// ИПарсерСек -- парсер секунд
type ИПарсерСек interface {
	ИПарсерПростой
	// Получ -- возвращает числовое значение секунд
	Получ() АСек
}
