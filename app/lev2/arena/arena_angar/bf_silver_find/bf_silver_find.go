// package bf_silver_find -- бизнес-функция поиск серебро бота
package bf_silver_find

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// СереброНайти -- ищет серебро бота
func СереброНайти(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	lstAngar := ангар.СписПолучить()
	var (
		strOut   string
		еслиЕсть bool
	)
	for _, strSilver := range lstAngar {
		if strings.Contains(strSilver, `<img title="Серебро" `) {
			strOut = strSilver
			еслиЕсть = true
			break
		}
	}
	Hassert(еслиЕсть, "СереброНайти()Ж строка серебра не найдена")
	// Выделить топливо
	lstSilver := strings.Split(strOut, `<img title="Серебро" alt="Серебро" src="/images/icons/silver.png?2"/> `)
	strSilver := lstSilver[1]
	iSilver, err := strconv.Atoi(strSilver)
	Hassert(err == nil, "СереброНайти(): серебро(%v) не число", iSilver)
	конт.Set("серебро", iSilver, "Серебро бота")
	ангар.Серебро().Уст(iSilver)
}
