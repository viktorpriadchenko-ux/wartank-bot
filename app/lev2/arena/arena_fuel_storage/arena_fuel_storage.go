// package arena_fuel_storage -- склад топлива
package arena_fuel_storage

import (
	"log"
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_fuel_storage/bf_fuel_build"
	"wartank/app/lev2/arena/arena_fuel_storage/bf_fuel_upgrade"
)

// АренаСкладТоплива -- склад топлива
type АренаБак struct {
	ИАренаСтроение
	конт ILocalCtx
	бак  ИСтатПарам
}

// НовТопливо -- возвращает новой склад топлива
func НовАренаБак(конт ILocalCtx) *АренаБак {
	сам := &АренаБак{
		конт: конт,
		бак:  lev1.НовСтатПарам("бак"),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Склад_топлива",
		СтрКонтроль_: `<title>Склад топлива</title>`,
		СтрУрл_:      "https://wartank.ru/fuelStore",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	конт.Set("бак", сам, "Топливная база")
	_ = ИАренаБак(сам)
	return сам
}

func (сам *АренаБак) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_fuel_build.БакПостроить(сам.конт)
	bf_fuel_upgrade.БакАпгрейд(сам.конт)
}

// Проверяет количество продукта в шахте
func (сам *АренаБак) количествоПолучить() {
	var (
		strOut      string
		еслиНайдено bool
	)
	lstMine := сам.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")
	/*
		<img class="ico vm" src="/images/icons/fuel.png?2" alt="Топливо" title="Топливо"> 720
	*/
	for _, strOut = range lstMine {
		if strings.Contains(strOut, `src=" alt="Топливо" title="Топливо"`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	_число := strings.TrimPrefix(strOut, `<img class="ico vm" src="/images/icons/fuel.png?2" alt="Топливо" title="Топливо"> `)
	iNum, err := strconv.Atoi(_число)
	if err != nil {
		log.Printf("АренаСкладТоплива.количествоПолучить(): кол-во топлива (%v) не число, err=\n\t%v\n", _число, err)
		return
	}
	сам.бак.Уст(iNum)
}

// Проверяет ускорение строительства
func (сам *АренаБак) ускорениеПровер() {
	списСтр := сам.Сеть().ВебВоркер().Получ("http://wartank.ru/buildings")
	// <span class="green2">Склад топлива -
	var (
		еслиНайти = false
		стр       string
	)
	for _, стр = range списСтр {
		if strings.Contains(стр, `<span class="green2">Склад топлива - `) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		return
	}
}

/*
// Пытается проапгрейдить топливный склад
func (сам *АренаБак) проапгрейдить() bool {
	time.Sleep(time.Millisecond * 1000)
	var (
		еслиНайти = false
		списСтр   []string
		стр       = ""
	)
	фнКупить := func() bool {
		defer time.Sleep(time.Millisecond * 1000)
		списСтр = сам.Сеть().ВебВоркер().Получ("https://wartank.ru/building-upgrade/FuelStorage")
		for _, стр = range списСтр {
			// <a class="simple-but border mb5" href="FuelStorage?5-1.ILinkListener-upgradeLink-link">
			if strings.Contains(стр, `ILinkListener-upgradeLink-link`) {
				еслиНайти = true
				break
			}
		}
		if !еслиНайти {
			return true
		}
		// Пробуем улучшить шахту
		_стр := strings.TrimPrefix(стр, "<a class=\"simple-but border mb5\" href=\"")
		_стр = strings.TrimSuffix(_стр, "\">")
		// https://wartank.ru/building-upgrade/FuelStorage?4-1.ILinkListener-upgradeLink-link
		// <a class="simple-but border mb5" href="FuelStorage?50-1.ILinkListener-upgradeLink-link">
		ссылка := "https://wartank.ru/building-upgrade/" + _стр
		списСтр = сам.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "ILinkListener-upgradeLink-link") {
				log.Printf("АренаСкладТоплива.проапгрейдить().фнКупить(): покупка склада топлива не прошла\n\tlink=%v\n\tстр=\n\t%v\n", ссылка, стр)
				return false // Покупка не оплачена
			}
		}
		log.Printf("+++++АренаСкладТоплива.проапгрейдить().фнКупить(): покупка склада топлива прошла\n")
		return true
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
		списСтр = сам.Сеть().ВебВоркер().Получ(ссылка)
		// Проверить, что постройка состоялась
		for _, стр := range списСтр {
			if strings.Contains(стр, "<title>Вы сделали слишком большую паузу</title>") {
				log.Printf("АренаСкладТоплива.проапгрейдить().фнПодтверждение(): подтверждение покупка склада топлива не прошла\n\tlink=%v\n\tстр=\n\t%v\n", ссылка, стр)
				return false // Покупка не оплачена
			}
		}
		log.Printf("+++++АренаСкладТоплива.проапгрейдить().фнПодтверждение(): подтверждение покупка склада топлива прошла\n")
		return true
	}

	фнКомплекс := func() {
		count := 5
		for count > 0 {
			if фнКупить() {
				if фнПодтверждение() {
					break
				}
			}
			count--
		}
	}
	фнКомплекс()
	return true
}
*/

// Бак -- возвращает размер бака
func (сам *АренаБак) Бак() ИСтатПарам {
	return сам.бак
}
