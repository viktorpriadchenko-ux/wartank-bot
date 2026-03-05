// package bf_fuel_build -- бизнес-функция строительства базы топлива
package bf_fuel_build

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БакПостроить -- постройка базы топлива
func БакПостроить(конт ILocalCtx) {
	бак := конт.Get("бак").Val().(ИАренаБак)
	if бак.Состояние().Получ() != cons.РежимНеСуществует {
		return
	}
	бакПостроить(конт)
}

func бакПостроить(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	бак := конт.Get("бак").Val().(ИАренаБак)
	списСтр := база.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")
	ссыльПостроить := "" // ссылка на постройку

	{ // Поиск ссылки на покупку
		// <td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/FuelStorage"><span><span>Построить</span></span></a></td>
		for _, стр := range списСтр {
			if strings.Contains(стр, `href="building-upgrade/FuelStorage">`) {
				ссыльПостроить = стр
				break
			}
		}
		if ссыльПостроить == "" {
			бак.Состояние().Уст(cons.РежимПостроено)
			return
		}
		// <td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/FuelStorage"><span><span>Построить</span></span></a></td>
		ссыльПостроить = strings.TrimPrefix(ссыльПостроить, `<td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="`)
		ссыльПостроить = strings.TrimSuffix(ссыльПостроить, `"><span><span>Построить</span></span></a></td>`)
		// https://wartank.ru/building-upgrade/FuelStorage
		ссыльПостроить = "http://wartank.ru/" + ссыльПостроить
		списСтр = база.Сеть().ВебВоркер().Получ(ссыльПостроить)
	}
	ссыльПодтвердить := "" // ссылка на улучшение здания

	{ // Выбор покупки
		// <a class="simple-but border mb5" href="FuelStorage?29-1.ILinkListener-upgradeLink-link">
		for _, стр := range списСтр {
			if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
				ссыльПодтвердить = стр
				break
			}
		}
		if ссыльПодтвердить == "" {
			бак.Состояние().Уст(cons.РежимПостроено)
			return
		}
		ссыльПодтвердить = strings.TrimPrefix(ссыльПодтвердить, `<a class="simple-but border mb5" href="`)
		ссыльПодтвердить = strings.TrimSuffix(ссыльПодтвердить, `">`)
		// https://wartank.ru/building-upgrade/FuelStorage?28-1.ILinkListener-upgradeLink-link
		ссыльПодтвердить = "http://wartank.ru/building-upgrade/" + ссыльПодтвердить
		списСтр = база.Сеть().ВебВоркер().Получ(ссыльПодтвердить)
	}
	ссыльДа := "" // подтверждение покупки
	{             // Подтверждение покупки
		// <a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../wicket/page?31-1.ILinkListener-confirmLink"><span><span>да, подтверждаю</span></span></a>
		for _, стр := range списСтр {
			if strings.Contains(стр, `confirmLink`) {
				ссыльДа = стр
				break
			}
		}
		if ссыльДа == "" {
			бак.Состояние().Уст(cons.РежимПостроено)
			return
		}
		ссыльДа = strings.TrimPrefix(ссыльДа, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../`)
		ссыльДа = strings.TrimSuffix(ссыльДа, `"><span><span>да, подтверждаю</span></span></a>`)
		// https://wartank.ru/wicket/page?52-1.ILinkListener-confirmLink
		ссыльДа = "http://wartank.ru/" + ссыльДа
		_ = база.Сеть().ВебВоркер().Получ(ссыльДа)
		бак.Состояние().Уст(cons.РежимПостроено)
	}
}
