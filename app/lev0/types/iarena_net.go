package types

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ИАренаСеть -- сетевые операции арены
type ИАренаСеть interface {
	ИБотСеть
	// Обновить -- обновляет список строк арены из сети
	Обновить()
	Get(strLink string) Result[[]string]
}
