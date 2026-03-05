// package bf_gold_by -- бизнес-функция покупки золота на рынке
package bf_gold_by

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// БанкПостроить -- покупка золота на рынке
func ЗолотоКупить(конт ILocalCtx) {
	рынок := конт.Get("рынок").Val().(ИАренаРынок)
	еслиПостроен := рынок.Состояние().Получ() == cons.РежимПостроено
	еслиОжидание := рынок.Состояние().Получ() == cons.РежимОжидание
	if !(еслиОжидание || еслиПостроен) {
		return
	}
	золотоКупить(конт)
}

func золотоКупить(конт ILocalCtx) {
	рынок := конт.Get("рынок").Val().(ИАренаРынок)
	var (
		ind         int
		еслиНайдено bool
		strOut      string
		lstMarket   = рынок.СписПолучить()
		strSilver   string
	)
	for ind, strOut = range lstMarket {
		if strings.Contains(strOut, `alt="Серебро" title="Серебро"> `) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Не найдена продажа золота за серебро
		return
	}
	lstSilver := strings.Split(strOut, `<img class="ico vm" src="/images/icons/silver.png?2" alt="Серебро" title="Серебро"> `)
	strSilver = lstSilver[1]
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	серебро := ангар.Серебро().Получ()
	еслиКупить := false
	switch strSilver {
	case "10", "50", "100", "500": // Допустимые суммы трат
		еслиКупить = true
	case "1000": // Если стоит тысяча серебра
		if серебро > 500_000 { // Если серебра больше полумиллиона -- покупаем
			еслиКупить = true
		}
	case "5000", "10000": // Если большая сумма -- можно купить и больше
		if серебро > 1_000_000 {
			еслиКупить = true
		}
	}
	if !еслиКупить {
		return
	}
	if ind < 15 {
		return // Недостаточно строк перед маркером
	}
	ind -= 15
	strOut = lstMarket[ind]
	lstLink := strings.Split(strOut, `<a class="simple-but border mb5" href="`)
	if len(lstLink) < 2 {
		return
	}
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить `)
	strLink = "https://wartank.ru/" + lstLink[0]
	res := рынок.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Market.buyGold(): при выполнении GET-команды на покупку золота, err=\n\t%v\n", err)
		return
	}
	lstMarket = res.Unwrap()
	for _, strOut = range lstMarket {
		if strings.Contains(strOut, `Ошибка на сервере. Сообщение админу уже отправлено.`) {
			// log._rintf("ERRO Market.buyGold(): при получении lstMarket, strHTML=%v, err=\nt%v\n", strOut, err)
			return
		}
	}
	рынок.СтрОбновить(lstMarket)
}
