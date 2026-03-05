// package bf_angar_fuel -- бизнес-функция поиска топлива в ангаре
package bf_angar_fuel

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ТопливоНайти -- возвращает топливо бота
func ТопливоНайти(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	if ангар.Состояние().Получ() == cons.РежимНеСуществует {
		ангар.Состояние().Уст(cons.РежимПостроено)
	}
	lstAngar := ангар.СписПолучить()
	var (
		strOut   string
		еслиЕсть bool
	)
	for _, strFuel := range lstAngar {
		if strings.Contains(strFuel, `<img title="Топливо" `) {
			strOut = strFuel
			еслиЕсть = true
			break
		}
	}
	Hassert(еслиЕсть, "ТопливоНайти(): нет строки топлива")
	// Выделить топливо
	lstFuel := strings.Split(strOut, `<img title="Топливо" alt="Топливо" src="/images/icons/fuel.png?2"/> `)
	// Здесь бывает ошибка (когда возвращена пустая строка)
	Hassert(len(lstFuel) == 2, "Обновить(): при поиске строки топлива, стр=\n\t%v\n", strOut)
	Hassert(lstFuel[1] != "", "Обновить(): пустое значение в строке топлива, стр=\n\t%v\n", strOut)
	strFuel := lstFuel[1]
	iFuel, err := strconv.Atoi(strFuel)
	Hassert(err == nil, "ТопливоНайти(): топливо(%v) не число", iFuel)
	ангар.Топливо().Уст(iFuel)
	if iFuel <= 30 { // Тратим топливо только если > 30
		return
	}
	конт.Set("топливо", iFuel, "Топливо бота")
}
