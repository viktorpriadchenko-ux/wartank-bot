package types

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ИСетьКлиент -- интерфейс к GET-запросу
type ИСетьКлиент interface {
	// Get -- теневая функция на блокировку
	Get(strLink string) Result[[]string]
}
