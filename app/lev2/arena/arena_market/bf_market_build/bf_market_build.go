// package bf_market_build -- бизнес-функция строительства рынка
package bf_market_build

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкПостроить -- постройка рынка
func РынокПостроить(конт ILocalCtx) {
	рынок := конт.Get("рынок").Val().(ИАренаРынок)
	if рынок.Состояние().Получ() == cons.РежимНеСуществует {
		рынокПостроить(конт)
	}
}

func рынокПостроить(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	рынок := конт.Get("рынок").Val().(ИАренаРынок)
	списСтр := база.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")
	ссыльПостроить := "" // ссылка на постройку

	{ // Поиск ссылки на покупку
		// <td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/Market"><span><span>Построить</span></span></a></td>
		for _, стр := range списСтр {
			if strings.Contains(стр, `href="building-upgrade/Market">`) {
				ссыльПостроить = стр
				break
			}
		}
		if ссыльПостроить == "" {
			рынок.Состояние().Уст(cons.РежимПостроено)
			return
		}
		// <td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/Market"><span><span>Построить</span></span></a></td>
		ссыльПостроить = strings.TrimPrefix(ссыльПостроить, `<td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="`)
		ссыльПостроить = strings.TrimSuffix(ссыльПостроить, `"><span><span>Построить</span></span></a></td>`)
		// https://wartank.ru/building-upgrade/Market
		ссыльПостроить = "http://wartank.ru/" + ссыльПостроить
		списСтр = база.Сеть().ВебВоркер().Получ(ссыльПостроить)
	}
	ссыльПодтвердить := "" // ссылка на улучшение здания

	{ // Выбор покупки
		// <a class="simple-but border mb5" href="Market?29-1.ILinkListener-upgradeLink-link">
		for _, стр := range списСтр {
			if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
				ссыльПодтвердить = стр
				break
			}
		}
		if ссыльПодтвердить == "" {
			рынок.Состояние().Уст(cons.РежимПостроено)
			return
		}
		ссыльПодтвердить = strings.TrimPrefix(ссыльПодтвердить, `<a class="simple-but border mb5" href="`)
		ссыльПодтвердить = strings.TrimSuffix(ссыльПодтвердить, `">`)
		// https://wartank.ru/building-upgrade/Market?28-1.ILinkListener-upgradeLink-link
		ссыльПодтвердить = "http://wartank.ru/building-upgrade/" + ссыльПодтвердить
		списСтр = база.Сеть().ВебВоркер().Получ(ссыльПодтвердить)
	}
	ссыльДа := "" // подтверждение покупки
	{             // Подтверждение покупки
		// <a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../wicket/page?195-1.ILinkListener-confirmLink"><span><span>да, подтверждаю</span></span></a>
		for _, стр := range списСтр {
			if strings.Contains(стр, `ILinkListener-confirmLink"><span><span>да, подтверждаю<`) {
				ссыльДа = стр
				break
			}
		}
		if ссыльДа == "" {
			рынок.Состояние().Уст(cons.РежимПостроено)
			return
		}
		ссыльДа = strings.TrimPrefix(ссыльДа, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../`)
		ссыльДа = strings.TrimSuffix(ссыльДа, `"><span><span>да, подтверждаю</span></span></a>`)
		// https://wartank.ru/wicket/page?195-1.ILinkListener-confirmLink
		ссыльДа = "http://wartank.ru/" + ссыльДа
		_ = база.Сеть().ВебВоркер().Получ(ссыльДа)
		рынок.Состояние().Уст(cons.РежимПостроено)
	}
}
