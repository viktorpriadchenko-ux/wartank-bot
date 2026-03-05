// package bf_glory_find -- поиск славы
package bf_glory_find

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// СлаваНайти -- ищет славу бота
func СлаваНайти(конт ILocalCtx) {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	// Найти строку с упоминанием оставшегося времени конвоя
	lstConvoy := конвой.СписПолучить()
	var (
		strGlory    string
		еслиНайдено bool
	)
	if len(lstConvoy) == 0 {
		конвой.Обновить()
		lstConvoy = конвой.СписПолучить()
	}
	for _, lastTime := range lstConvoy {
		if strings.Contains(lastTime, `alt="Слава" title="Слава"> `) {
			strGlory = lastTime
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Может не быть, если началась атака
		return
	}
	// Ищем количество славы
	// <img class="ico vm" src="/images/icons/glory.png?2" alt="Слава" title="Слава"> 167
	strGlory = strings.TrimPrefix(strGlory, `<img class="ico vm" src="/images/icons/glory.png?2" alt="Слава" title="Слава"> `)
	iGlory, err := strconv.Atoi(strGlory)
	Hassert(err == nil, "СлаваНайти(): слава(%v) не число, err=\n\t%v\n", strGlory, err)
	танкСтат := конт.Get("танкСтат").Val().(ИТанкСтат)
	танкСтат.Слава().Уст(iGlory)
}
