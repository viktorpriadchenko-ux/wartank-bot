// package bf_mine_accelerate -- бизнес-функция ускорения строительства или апгрейда
package bf_mine_accelerate

import (
	"fmt"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ШахтаУскорить -- пробует ускорить строительство шахты или апгрейда
func ШахтаУскорить(конт ILocalCtx) {
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	еслиАпгрейд := шахта.Состояние().Получ() == cons.РежимАпгрейд
	еслиПлатный := шахта.Состояние().Получ() == cons.РежимАпгрейдПлатный
	if !(еслиАпгрейд || еслиПлатный) {
		return
	}
	шахтаАпгрейд(конт)
}

func шахтаАпгрейд(конт ILocalCtx) {
	база := конт.Get("база").Val().(ИАренаБаза)
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	var (
		еслиНайти = false
		списСтр   []string
		стр       = ""
	)
	фнКупитьАпгрейд := func() (bool, error) {
		списСтр = база.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/Mine")
		for _, стр = range списСтр {
			// <a class="simple-but border mb5" href="Mine?5-1.ILinkListener-upgradeLink-link">
			if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			return false, nil
		}
		// Пробуем улучшить шахту
		_стр := strings.TrimPrefix(стр, "<a class=\"simple-but border mb5\" href=\"")
		_стр = strings.TrimSuffix(_стр, "\">")
		// https://wartank.ru/building-upgrade/Mine?4-1.ILinkListener-upgradeLink-link
		// <a class="simple-but border mb5" href="FuelStorage?50-1.ILinkListener-upgradeLink-link">
		ссылка := "https://wartank.ru/building-upgrade/" + _стр
		списСтр = база.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "ILinkListener-upgradeLink-link") {
				return false, fmt.Errorf("покупка шахты не прошла") // Покупка не оплачена
			}
		}
		return true, nil
	}

	фнПодтверждение := func() bool {
		for _, стр = range списСтр {
			// <a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="../wicket/page?7-1.ILinkListener-confirmLink"><span><span>да, подтверждаю</span></span></a>
			if strings.Contains(стр, `ILinkListener-confirmLink`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			return true
		}
		// Пробуем построить шахту
		_стр := strings.TrimPrefix(стр, `<a class="simple-but border w50 mXa mb10" w:id="confirmLink" href="..`)
		_стр = strings.TrimSuffix(_стр, `"><span><span>да, подтверждаю</span></span></a>`)
		// https://wartank.ru/wicket/page?6-1.ILinkListener-confirmLink
		ссылка := "https://wartank.ru" + _стр
		списСтр = база.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "<title>Вы сделали слишком большую паузу</title>") {
				return false // Покупка не оплачена
			}
		}
		return true
	}

	фнКомплекс := func() {
		еслиОк, ош := фнКупитьАпгрейд()
		switch {
		case ош == nil && еслиОк: // покупка апгрейда шахты прошла
			if фнПодтверждение() {
				шахта.Состояние().Уст(cons.РежимАпгрейдПлатный)
				return
			}
		case ош == nil && !еслиОк: // покупка шахты не нужна
			шахта.Состояние().Уст(cons.РежимОжидание)
			return
		case ош != nil: // ошибка при работе с сетью
			шахта.Состояние().Уст(cons.РежимАпгрейд)
			return
		}
	}
	фнКомплекс()
}
