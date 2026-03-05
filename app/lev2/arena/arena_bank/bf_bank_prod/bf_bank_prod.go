// package bf_bank_prod -- бизнес-функция производить серебро
package bf_bank_prod

import (
	"log"
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// СереброПроизводить -- заставляет банк производить серебро
func СереброПроизводить(конт ILocalCtx) {
	банк := конт.Get("банк").Val().(ИАренаБанк)
	еслиОжидание := банк.Состояние().Получ() == cons.РежимОжидание
	еслиПостроен := банк.Состояние().Получ() == cons.РежимПостроено
	if !(еслиОжидание || еслиПостроен) {
		return
	}
	получитьВсеРежимы(конт)
	сделатьСеребро(конт)
}

// Получает все режимы банка
func получитьВсеРежимы(конт ILocalCtx) {
	банк := конт.Get("банк").Val().(ИАренаБанк)

	var (
		lstBank  = банк.СписПолучить()
		ind      int
		strMode  string
		strMode1 string
		strMode2 string
	)
	if len(lstBank) == 0 {
		банк.Обновить()
		lstBank = банк.СписПолучить()
	}
	{ // Получить первый режим
		for ind, strMode = range lstBank {
			if strings.Contains(strMode, `Кол-во: <span class="green2">`) {
				strMode1 = strMode
				break
			}
		}
		lstMode := strings.Split(strMode1, `Кол-во: <span class="green2">`)
		strMode1 = lstMode[1]
		lstMode = strings.Split(strMode1, `</span><br/>`)
		strMode1 = lstMode[0]
		iNum1, err := strconv.Atoi(strMode1)
		Hassert(err == nil, "получитьВсеРежимы(): ошибка в режиме-1 банка, ош=\n\t%v", err)
		банк.РежимРаботы1().Серебро().Уст(iNum1)
		// Установить время производства
		strTime1 := lstBank[ind+2]
		банк.РежимРаботы1().ВремяСделатьУст(strTime1)
		if iNum1 <= 2 { // Если банк слишком слабый
			return
		}
	}
	{ // Получить второй режим
		for _ind := ind + 2; _ind < len(lstBank); _ind++ {
			strMode := lstBank[_ind]
			if strings.Contains(strMode, `Кол-во: <span class="green2">`) {
				strMode2 = strMode
				ind = _ind
				break
			}
		}
		lstMode := strings.Split(strMode2, `Кол-во: <span class="green2">`)
		strMode2 = lstMode[1]
		lstMode = strings.Split(strMode2, `</span><br/>`)
		strMode2 = lstMode[0]
		iNum2, err := strconv.Atoi(strMode2)
		Hassert(err == nil, "получитьВсеРежимы(): ошибка в режиме-2 банка, ош=\n\t%v", err)
		банк.РежимРаботы2().Серебро().Уст(iNum2)
		// Установить время производства
		strTime2 := lstBank[ind+2]
		банк.РежимРаботы2().ВремяСделатьУст(strTime2)
	}
}

// Запускает в производство серебро
func сделатьСеребро(конт ILocalCtx) {
	банк := конт.Get("банк").Val().(ИАренаБанк)
	var (
		lstBank = банк.СписПолучить()
		ind     int
		strOut  string
		strLink string
	)
	time1 := банк.РежимРаботы1().ВремяСделать()
	time2 := банк.РежимРаботы2().ВремяСделать()
	if time1 > time2 {
		time1 = time2
	}
	fnRetry := func() bool {
		for ind, strOut = range lstBank {
			if strings.Contains(strOut, time1) {
				ind += 7
				strLink = lstBank[ind]
				break
			}
		}
		return strings.Contains(strLink, `>Начать производство</span>`)
	}
	for fnRetry() {
		банк.Обновить()
		lstBank = банк.СписПолучить()
		lstLink := strings.Split(strLink, `<a class="simple-but border" href="`)
		strLink = lstLink[1]
		lstLink = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
		strLink = "https://wartank.ru/production/" + lstLink[0]
		log.Printf("сделатьСеребро(): link=%v", strLink)
		lstBank = банк.Сеть().ВебВоркер().Получ(strLink)
		log.Printf("сделатьСеребро(): link=%v", strLink)
	}
	банк.ПродуктВремяСейчас().Set(time1)
}
