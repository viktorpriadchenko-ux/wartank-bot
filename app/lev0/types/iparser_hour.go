package types

import "wartank/app/lev0/alias"

// ИПарсерЧас --парсер часов
type ИПарсерЧас interface {
	ИПарсерПростой
	// Получ -- возвращает числовое значение часа
	Получ() alias.АЧас
}
