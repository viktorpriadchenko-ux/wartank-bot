package types

import "wartank/app/lev0/alias"

// ИПарсерМин -- парсер минут
type ИПарсерМин interface {
	ИПарсерПростой
	// Получ -- возвращает числовое значение минут
	Получ() alias.АМин
}
