// package bf_polygon_upgrade_fast -- бизнес-функция бесплатного апгрейда полигона
package bf_polygon_upgrade_fast

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкАпгрейд -- ускоряет повышение уровня полигона бесплатно
func ПолигонАпгрейдБесплатно(конт ILocalCtx) {
	полигон := конт.Get("полигон").Val().(ИАренаПолигон)
	еслиПлатно := полигон.Состояние().Получ() == cons.РежимАпгрейдПлатный
	еслиПостроено := полигон.Состояние().Получ() == cons.РежимПостроено
	if !(еслиПлатно || еслиПостроено) {
		return
	}
	полигонАпгрейдБесплатно(конт)
}

// https://wartank.ru/buildings?16-1.ILinkListener-buildings-2-building-rootBlock-actionPanel-freeBoostLink
func полигонАпгрейдБесплатно(конт ILocalCtx) {
	полигон := конт.Get("полигон").Val().(ИАренаПолигон)
	var (
		стрСсылка   = ""
		еслиНайдено = false
	)
	списБанк := полигон.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")
	// <td style="width:50%;padding-left:1px;"><a class="simple-but border" href="buildings?61-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-freeBoostLink"><span><span>Ускорение</span></span></a>
	for _, стрСсылка = range списБанк {
		if strings.Contains(стрСсылка, `building-rootBlock-actionPanel-freeBoostLink"><span><span>Ускорение<`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	lstLink := strings.Split(стрСсылка, `<td style="width:50%;padding-left:1px;"><a class="simple-but border" href="`)
	if len(lstLink) < 2 {
		return
	}
	_ссылка := lstLink[1]
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Ускорение</span></span></a>`)
	// https://wartank.ru/buildings?61-1.ILinkListener-buildings-0-building-rootBlock-actionPanel-freeBoostLink
	ссылка := "https://wartank.ru/" + _ссылка
	_ = полигон.Сеть().ВебВоркер().Получ(ссылка)
	if полигон.Состояние().Получ() == cons.РежимПостроено {
		полигон.Состояние().Уст(cons.РежимАпгрейдПлатный)
	}
}
