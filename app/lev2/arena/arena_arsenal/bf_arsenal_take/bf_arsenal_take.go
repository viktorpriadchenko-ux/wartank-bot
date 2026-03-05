// package bf_arsenal_take -- бизнес-функция забрать из арсенала
package bf_arsenal_take

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкЗабрать -- забрать из арсенала
func АрсеналЗабрать(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	if арсенал.Состояние().Получ() != cons.РежимЗабрать {
		return
	}
	арсеналЗабрать(конт)
}

func арсеналЗабрать(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	var (
		strOut      string
		еслиНайдено bool
		списСтр     []string
	)
	списСтр = база.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")

	// <a class="simple-but border" href="buildings?35-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-takeProductionLink"><span><span>Забрать</span></span></a>
	for _, strOut = range списСтр {
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
	сост := арсенал.Состояние().Получ()
	if сост == cons.РежимНеСуществует {
		арсенал.Состояние().Уст(cons.РежимПостроено)
	}
	if арсенал.Состояние().Получ() == cons.РежимРабота {
		арсенал.Состояние().Уст(cons.РежимЗабрать)
	}
	арсенал.Состояние().Уст(cons.РежимОжидание)
}
