// package bf_bank_take -- бизнес-функция забрать из банка
package bf_bank_take

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкЗабрать -- забрать из банка
func БанкЗабрать(конт ILocalCtx) {
	банк := конт.Get("банк").Val().(ИАренаБанк)
	еслиПостроено := банк.Состояние().Получ() == cons.РежимПостроено
	еслиЗабрать := банк.Состояние().Получ() == cons.РежимЗабрать
	if !(еслиПостроено || еслиЗабрать) {
		return
	}
	банкЗабрать(конт)
}

func банкЗабрать(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	банк := конт.Get("банк").Val().(ИАренаБанк)
	strOut := ""
	списБанк := база.СписПолучить()
	if len(списБанк) == 0 {
		база.ОбновитьПринуд()
		списБанк = база.СписПолучить()
	}
	адр := 0
	фнНайтиЗабрать := func() bool {
		// 119 - Производит серебро<br/>
		// 136 - <a class="simple-but border" href="buildings?91-1.ILinkListener-buildings-1-building-rootBlock-actionPanel-takeProductionLink"><span><span>Забрать</span></span></a>
		for адр, strOut = range списБанк {
			if strings.Contains(strOut, `Производит серебро<br/>`) {
				strOut = списБанк[адр+17]
				return true
			}
		}
		return false
	}
	if фнНайтиЗабрать() {
		if !strings.Contains(strOut, `Забрать`) {
			return
		}
		_ссылка := strings.TrimPrefix(strOut, `<a class="simple-but border" href="`)
		_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Забрать</span></span></a>`)
		// За время опроса/поиска -- ситуация могла измениться
		ссылка := "https://wartank.ru/" + _ссылка
		// http://wartank.ru/buildings?5-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink
		списБанк = база.Сеть().ВебВоркер().Получ(ссылка)
		for фнНайтиЗабрать() {
			_ссылка := strings.TrimPrefix(strOut, `<a class="simple-but border" href="`)
			_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Забрать</span></span></a>`)
			ссылка := "https://wartank.ru/" + _ссылка
			// http://wartank.ru/buildings?5-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink
			списБанк = база.Сеть().ВебВоркер().Получ(ссылка)
		}
		еслиПостроено := банк.Состояние().Получ() == cons.РежимПостроено
		еслиРабота := банк.Состояние().Получ() == cons.РежимРабота
		if еслиПостроено || еслиРабота {
			банк.Состояние().Уст(cons.РежимЗабрать)
			банк.Состояние().Уст(cons.РежимОжидание)
		}

	}

}
