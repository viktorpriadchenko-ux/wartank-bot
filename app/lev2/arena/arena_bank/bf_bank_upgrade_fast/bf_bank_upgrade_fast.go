// package bf_bank_upgrade_fast -- бизнес-функция бесплатного апгрейда банка
package bf_bank_upgrade_fast

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкАпгрейд -- ускоряет повышение уровня банка бесплатно
func БанкАпгрейдБесплатно(конт ILocalCtx) {
	банк := конт.Get("банк").Val().(ИАренаБанк)
	еслиПлатно := банк.Состояние().Получ() == cons.РежимАпгрейдПлатный
	еслиПостроено := банк.Состояние().Получ() == cons.РежимПостроено
	if !(еслиПлатно || еслиПостроено) {
		return
	}
	банкАпгрейдБесплатно(конт)
}

// https://wartank.ru/buildings?16-1.ILinkListener-buildings-2-building-rootBlock-actionPanel-freeBoostLink
func банкАпгрейдБесплатно(конт ILocalCtx) {
	банк := конт.Get("банк").Val().(ИАренаБанк)
	var (
		стрСсылка   = ""
		еслиНайдено = false
	)
	списБанк := банк.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")
	// <td style="width:50%;padding-left:1px;"><a class="simple-but border" href="buildings?1-1.ILinkListener-buildings-1-building-rootBlock-actionPanel-freeBoostLink"><span><span>Ускорение</span></span></a>
	for _, стрСсылка = range списБанк {
		if strings.Contains(стрСсылка, `.ILinkListener-buildings-1-building-rootBlock-actionPanel-freeBoostLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	_ссылка := strings.TrimPrefix(стрСсылка, `<td style="width:50%;padding-left:1px;"><a class="simple-but border" href="`)
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Ускорение</span></span></a>`)
	ссылка := "https://wartank.ru/" + _ссылка
	_ = банк.Сеть().ВебВоркер().Получ(ссылка)
	if банк.Состояние().Получ() == cons.РежимПостроено {
		банк.Состояние().Уст(cons.РежимАпгрейдПлатный)
	}
}
