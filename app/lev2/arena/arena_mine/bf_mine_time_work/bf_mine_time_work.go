// package bf_mine_time -- вычисляет оставшееся время производства
package bf_mine_time_work

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ШахтаРаботаВремя -- выясняет оставшееся время
func ШахтаРаботаВремя(конт ILocalCtx) {
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	if шахта.Состояние().Получ() != cons.РежимРабота {
		return
	}
	база := конт.Get("база").Val().(ИАренаБаза)
	lstMine := база.СписПолучить()
	if len(lstMine) == 0 {
		база.ОбновитьПринуд()
		lstMine = база.СписПолучить()
	}
	var (
		ind    int
		str    string
		isFind bool
	)
	// <span class="green2">Шахта - 1</span><br/>
	for ind, str = range lstMine {
		if strings.Contains(str, `<span class="green2">Шахта - `) {
			isFind = true
			break
		}
	}
	Hassert(isFind, "ШахтаВремя(): строка времени не найдена")

	// <td><div class="value-block lh1"><span><span>00:00:34</span></span></div></td>
	strTime := lstMine[ind+11]
	if !strings.Contains(strTime, ":") { // Уже время производства закончилось
		шахта.Состояние().Уст(cons.РежимЗабрать)
		return
	}
	// <td><div class="value-block lh1"><span><span>00:19:53</span></span></div></td>
	strTime = strings.TrimPrefix(strTime, `<td><div class="value-block lh1"><span><span>`)
	strTime = strings.TrimSuffix(strTime, `</span></span></div></td>`)
	// Здесь уже может выйти время работы, нужна проверка на контрольную строку
	if strTime == `<td style="width:50%;padding-right:1px;">` {
		return
	}
	шахта.ОбратВремяУст(АВремя(strTime))
	шахта.ВебЛог().Добавить("Шахта.количествоПолучить(): время=%q\n", strTime)
}
