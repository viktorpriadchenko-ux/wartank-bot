// package arena_market -- объект рынка
package arena_market

import (
	"log"
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_market/bf_gold_by"
	"wartank/app/lev2/arena/arena_market/bf_market_build"
	"wartank/app/lev2/arena/arena_market/bf_market_upgrade"
)

// АренаРынок -- объект рынка
type АренаРынок struct {
	ИАренаСтроение
	конт ILocalCtx
}

// НовРынок -- возвращает новый рынок
func НовРынок(конт ILocalCtx) ИАренаРынок {
	сам := &АренаРынок{
		конт: конт,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Рынок",
		СтрКонтроль_: `<title>Рынок</title>`,
		СтрУрл_:      "https://wartank.ru/market",
	}
	конт.Set("рынок", сам, "Рынок бота")
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	return сам
}

// Пуск -- запускает всю работу рынка в отдельном потоке
func (сам *АренаРынок) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_market_build.РынокПостроить(сам.конт)
	bf_market_upgrade.РынокАпгрейд(сам.конт)
	bf_gold_by.ЗолотоКупить(сам.конт)
}

// Проверяет ускорение строительства
func (сам *АренаРынок) ускорениеПровер() {
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

// Обновляет текущий уровень рынка (может быть не построен)
func (сам *АренаРынок) уровеньОбновить() bool {
	списСтр := сам.Сеть().ВебВоркер().Получ("http://wartank.ru/buildings")
	// <span class="green2">Рынок -
	var (
		еслиНайти = false
		стр       = ""
	)
	for _, стр = range списСтр {
		if strings.Contains(стр, `<span class="green2">Рынок -`) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		return false
	}
	// <span class="green2">Рынок - 0</span><br/>
	_стр := strings.TrimPrefix(стр, `<span class="green2">Рынок - `)
	_стр = strings.TrimSuffix(_стр, `</span><br/>`)
	иУровень, ош := strconv.Atoi(_стр)
	if ош != nil {
		log.Printf("Рынок.уровеньОбновить(): строка уровня сбойная, стр=%q, ош=\n\t%v\n", стр, ош)
		return false
	}
	сам.Уровень().Уст(иУровень)
	return true
}

// Проверяет  время ожидания рынка
func (сам *АренаРынок) проверОжидание() {
	var (
		strOut      string
		еслиНайдено bool
	)
	// countDown := сам.CountDown().Get()
	фнЕслиСеребро := func() bool { // Найти счётчик цены серебра
		сам.Обновить()
		еслиНайдено := false
		lstMarket := сам.СписПолучить()
		for _, strOut = range lstMarket {
			if strings.Contains(strOut, `alt="Серебро" title="Серебро"> `) {
				еслиНайдено = true
				break
			}
		}
		if еслиНайдено {
			lstSilver := strings.Split(strOut, `<img class="ico vm" src="/images/icons/silver.png?2" alt="Серебро" title="Серебро"> `)
			strSilver := lstSilver[1]
			switch strSilver {
			case "10", "50", "100", "500":
				return true
			default:
				серебро := сам.конт.Get("серебро").Val().(int)
				if серебро > 1_000_000 {
					return true
				}
				return false
			}
		}
		return false
	}

	fnGetCountDown := func() { // Искать счётчик времени
		lstMarket := сам.СписПолучить()
		// Найти счётчик времени
		for _, strOut = range lstMarket {
			if strings.Contains(strOut, `Минимальная цена через `) {
				еслиНайдено = true
				break
			}
		}
		if !еслиНайдено {
			return // Минимальная цена
		}
		lstTime := strings.Split(strOut, `Минимальная цена через `)
		strTime := lstTime[1]
		сам.ОбратВремяУст(АВремя(strTime))
	}
	if фнЕслиСеребро() {
		return
	}
	fnGetCountDown()
}
