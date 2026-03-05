// package lev4 -- сборочный слой
package lev4

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev4/mod_serv"
)

// НовМодСервер -- возвращает новый модуль сервера
func НовМодСервер() IKernelModule {
	сервер := mod_serv.НовМодСервер()
	return сервер
}
