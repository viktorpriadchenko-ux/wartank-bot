// package bf_market_upgrade -- бизнес-функция апгрейда рынка
package bf_market_upgrade

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкАпгрейд -- повышает уровень рынка
func РынокАпгрейд(конт ILocalCtx) {
	рынок := конт.Get("рынок").Val().(ИАренаРынок)
	еслиПостроено := рынок.Состояние().Получ() == cons.РежимПостроено
	еслиОжидание := рынок.Состояние().Получ() == cons.РежимОжидание
	if !(еслиПостроено || еслиОжидание) {
		return
	}
	рынокАпгрейд(конт)
}

func рынокАпгрейд(конт ILocalCtx) {
	рынок := конт.Get("рынок").Val().(ИАренаРынок)
	var (
		еслиНайти = false
		списСтр   []string
		стр       = ""
	)
	фнКупить := func() bool {
		списСтр = рынок.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Market")
		for _, стр = range списСтр {
			// <a class="simple-but border mb5" href="Market?5-1.ILinkListener-upgradeLink-link">
			if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			return true
		}
		// Пробуем улучшить здание
		_стр := strings.TrimPrefix(стр, "<a class=\"simple-but border mb5\" href=\"")
		_стр = strings.TrimSuffix(_стр, "\">")
		// https://wartank.ru/building-upgrade/Market?4-1.ILinkListener-upgradeLink-link
		// <a class="simple-but border mb5" href="Market?50-1.ILinkListener-upgradeLink-link">
		ссылка := "https://wartank.ru/building-upgrade/" + _стр
		списСтр = рынок.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что покупка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "ILinkListener-upgradeLink-link") {
				return false // Покупка не оплачена
			}
		}
		return true
	}

	фнПодтверждение := func() {
		for _, стр = range списСтр {
			// <a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../wicket/page?7-1.ILinkListener-confirmLink"><span><span>да, подтверждаю</span></span></a>
			if strings.Contains(стр, `ILinkListener-confirmLink`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			return
		}
		// Пробуем оплатить апгрейд
		_стр := strings.TrimPrefix(стр, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="..`)
		_стр = strings.TrimSuffix(_стр, `"><span><span>да, подтверждаю</span></span></a>`)
		// https://wartank.ru/wicket/page?6-1.ILinkListener-confirmLink
		ссылка := "https://wartank.ru" + _стр
		списСтр = рынок.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что оплата состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "<title>Вы сделали слишком большую паузу</title>") {
				return // Покупка не оплачена
			}
		}
		рынок.Состояние().Уст(cons.РежимАпгрейдПлатный)
	}

	if !фнКупить() {
		return
	}
	фнПодтверждение()
}
