// package product -- объект продукта для производства
package product

import (
	. "wartank/app/lev0/types"
	"wartank/app/lev1/product/parser_time"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// Продукт -- объект продукта для производства
type Продукт struct {
	имя   ISafeString              // имя продукта
	кол   ISafeInt                 // Количество продукта
	время *parser_time.ПарсерВремя // Время для производства продукта
}

// НовПродукт -- возвращает новый объект продукта
func НовПродукт() *Продукт {
	сам := &Продукт{
		имя:   NewSafeString(),
		кол:   NewSafeInt(),
		время: parser_time.НовПарсерВремя(),
	}
	return сам
}

// Имя -- возвращает название продукта
func (сам *Продукт) Имя() ISafeString {
	return сам.имя
}

// Кол -- возвращает количество продукта
func (сам *Продукт) Кол() ISafeInt {
	return сам.кол
}

// Время -- возвращает время производства
func (сам *Продукт) Время() ИПарсерВремя {
	return сам.время
}
