// package bf_polygon_upgrade -- бизнес-функция апгрейда полигона
package bf_polygon_upgrade

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ПолигонАпгрейд -- повышает уровень полигона
func ПолигонАпгрейд(конт ILocalCtx) {
	полигон := конт.Get("полигон").Val().(ИАренаПолигон)
	еслиПостроено := полигон.Состояние().Получ() == cons.РежимПостроено
	еслиОжидание := полигон.Состояние().Получ() == cons.РежимОжидание
	if !(еслиПостроено || еслиОжидание) {
		return
	}
	полигонАпгрейд(конт)
}

func полигонАпгрейд(конт ILocalCtx) {
	полигон := конт.Get("полигон").Val().(ИАренаПолигон)
	var (
		еслиНайти = false
		списСтр   []string
		стр       = ""
	)
	фнКупить := func() bool {
		списСтр = полигон.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Polygon")
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
		// Пробуем улучшить полигон
		_стр := strings.TrimPrefix(стр, "<a class=\"simple-but border mb5\" href=\"")
		_стр = strings.TrimSuffix(_стр, "\">")
		// https://wartank.ru/building-upgrade/Polygon?4-1.ILinkListener-upgradeLink-link
		// <a class="simple-but border mb5" href="Polygon?50-1.ILinkListener-upgradeLink-link">
		ссылка := "https://wartank.ru/building-upgrade/" + _стр
		списСтр = полигон.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
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
		// Пробуем подтвердить оплату
		_стр := strings.TrimPrefix(стр, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="..`)
		_стр = strings.TrimSuffix(_стр, `"><span><span>да, подтверждаю</span></span></a>`)
		// https://wartank.ru/wicket/page?6-1.ILinkListener-confirmLink
		ссылка := "https://wartank.ru" + _стр
		списСтр = полигон.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что оплата состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "<title>Вы сделали слишком большую паузу</title>") {
				return // Покупка не оплачена
			}
		}
		полигон.Состояние().Уст(cons.РежимАпгрейдПлатный)
	}

	if !фнКупить() {
		return
	}
	фнПодтверждение()
}
