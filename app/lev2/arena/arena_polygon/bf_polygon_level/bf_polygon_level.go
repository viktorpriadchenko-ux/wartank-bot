// package bf_polygon_level -- бизнес-функция, ищет уровень полигона
package bf_polygon_level

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ПолигонУровень -- бизнес-функция, ищет уровень полигона
func ПолигонУровень(конт ILocalCtx) {
	полигон := конт.Get("полигон").Val().(ИАренаПолигон)
	if полигон.Состояние().Получ() == cons.РежимНеСуществует {
		полигон.Уровень().Уст(0)
		return
	}
	var (
		стрВых      = ""
		еслиНайдено bool
	)
	база := конт.Get("база").Val().(ИАренаБаза)
	lstBase := база.СписПолучить()
	if len(lstBase) == 0 {
		база.ОбновитьПринуд()
		lstBase = база.СписПолучить()
	}
	// <span class="green2">Полигон - 5</span><br/>
	for _, стрВых = range lstBase {
		if strings.Contains(стрВых, `<span class="green2">Полигон - `) {
			еслиНайдено = true
			break
		}
	}
	Hassert(еслиНайдено, "ПолигонУровень(): не найдена контрольная строка")
	стрУровень := strings.TrimPrefix(стрВых, `<span class="green2">Полигон - `)
	стрУровень = strings.TrimSuffix(стрУровень, `</span><br/>`)
	цУров, ош := strconv.Atoi(стрУровень)
	Hassert(ош == nil, "ПолигонУровень(): уровень(%v) не число, ош=\n\t", стрУровень, ош)
	полигон.Уровень().Уст(цУров)
}
