// package bf_gold_find -- бизнес-функция поиска золота
package bf_gold_find

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// ЗолотоНайти -- ищет золото бота
func ЗолотоНайти(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	var (
		lstAngar    = ангар.СписПолучить()
		strOut      string
		еслиНайдено bool
	)
	if len(lstAngar) == 0 {
		ангар.Обновить()
		lstAngar = ангар.СписПолучить()
	}
	for _, strGold := range lstAngar {
		if strings.Contains(strGold, `<img title="Золото" `) {
			strOut = strGold
			еслиНайдено = true
			break
		}
	}
	Hassert(еслиНайдено, "ЗолотоНайти(): не найдена строка золота")
	// Выделить топливо
	lstGold := strings.Split(strOut, `<img title="Золото" alt="Золото" src="/images/icons/gold.png?2"/> `)
	strGold := lstGold[1]
	iGold, err := strconv.Atoi(strGold)
	Hassert(err == nil, "ЗолотоНайти(): iGold(%v) не число, ош=%v", iGold, err)
	ангар.Золото().Уст(iGold)
	конт.Set("золото", iGold, "Золото бота")
}
