// package bf_mine_build -- бизнес-функция постройки шахты
package bf_mine_build

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ШахтаПостроить -- постройка шахты
func ШахтаПостроить(конт ILocalCtx) {
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	if шахта.Состояние().Получ() == cons.РежимНеСуществует {
		шахтаПостроить(конт)
	}
}

func шахтаПостроить(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	// <td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/Mine"><span><span>Построить</span></span></a></td>
	var (
		еслиНайти = false
		стр       = ""
	)
	списСтр := база.Сеть().ВебВоркер().Получ("http://wartank.ru/buildings")
	for _, стр = range списСтр {
		if strings.Contains(стр, `href="building-upgrade/Mine"><span><span>Построить</span></span>`) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		шахта.Состояние().Уст(cons.РежимПостроено)
		return
	}
	// Пробуем построить шахту
	_стр := strings.TrimPrefix(стр, `<td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="`)
	_стр = strings.TrimSuffix(_стр, `"><span><span>Построить</span></span></a></td>`)
	ссылка := "https://wartank.ru/" + _стр
	списСтр = база.Сеть().ВебВоркер().Получ(ссылка)
	еслиНайти = false
	// "<a class=\"simple-but border mb5\" href=\"Mine?14-1.ILinkListener-upgradeLink-link\">"
	for _, стр = range списСтр {
		if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		шахта.Состояние().Уст(cons.РежимПостроено)
		return
	}
	_стр = strings.TrimPrefix(стр, "<a class=\"simple-but border mb5\" href=\"")
	_стр = strings.TrimSuffix(_стр, "\">")
	// http://wartank.ru/building-upgrade/Mine?16-1.ILinkListener-upgradeLink-link
	ссылка = "https://wartank.ru/building-upgrade/" + _стр
	_ = база.Сеть().ВебВоркер().Получ(ссылка)
	шахта.Состояние().Уст(cons.РежимПостроено)
}
