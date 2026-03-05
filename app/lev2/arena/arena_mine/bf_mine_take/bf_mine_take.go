// package bf_mine_take -- бизнес-функция забрать из шахты
package bf_mine_take

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ШахтаЗабрать -- забрать из шахты
func ШахтаЗабрать(конт ILocalCtx) {
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	if шахта.Состояние().Получ() != cons.РежимЗабрать {
		return
	}
	шахтаЗабрать(конт)
}

func шахтаЗабрать(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	var (
		strOut      string
		еслиНайдено bool
		списШахта   []string
	)
	списШахта = база.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")

	// <a class="simple-but border" href="buildings?35-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink"><span><span>Забрать</span></span></a>
	for _, strOut = range списШахта {
		if strings.Contains(strOut, `.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	_ссылка := strings.TrimPrefix(strOut, `<a class="simple-but border" href="`)
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Забрать</span></span></a>`)
	ссылка := "https://wartank.ru/" + _ссылка
	// http://wartank.ru/buildings?5-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink
	_ = база.Сеть().ВебВоркер().Получ(ссылка)
	сост := шахта.Состояние().Получ()
	if сост == cons.РежимНеСуществует {
		шахта.Состояние().Уст(cons.РежимПостроено)
	}
	if шахта.Состояние().Получ() == cons.РежимРабота {
		шахта.Состояние().Уст(cons.РежимЗабрать)
	}
	шахта.Состояние().Уст(cons.РежимОжидание)
}
