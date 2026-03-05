// package lev1 -- слой сущностей, компонентов
package lev1

import (
	. "wartank/app/lev0/types"
	"wartank/app/lev1/stat_param"
)

// НовСтатПарам -- возвращает новый ИСтатПарам
func НовСтатПарам(имя string) ИСтатПарам {
	return stat_param.НовСтатПарам1(имя)
}
