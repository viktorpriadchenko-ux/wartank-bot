// package bf_masters_wait -- ожидание начала битвы мастеров
package bf_masters_wait

import (
	"fmt"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// МастераОжидать -- ожидает начала битвы мастеров
func МастераОжидать(конт ILocalCtx) {
	битва := конт.Get("pvp").Val().(ИАренаСтроение)
	/*
		Целевые состояния: РежимАпгрейдПлатный (зарегистрированы, ждём старта)
		или РежимПостроено (свежий запуск, попробуем найти таймер).
	*/
	еслиПостроено := битва.Состояние().Получ() == cons.РежимПостроено
	еслиАпгрейд := битва.Состояние().Получ() == cons.РежимАпгрейдПлатный
	if !(еслиАпгрейд || еслиПостроено) {
		return
	}
	ждать(конт)
}

// ждать -- устанавливает таймер обратного отсчёта до начала битвы мастеров
func ждать(конт ILocalCtx) {
	битва := конт.Get("pvp").Val().(ИАренаСтроение)
	var (
		strOut      string
		lstBattle   = битва.СписПолучить()
		еслиНайдено bool
	)
	// <span>до начала 00:46:42 (2 993 заявок)</span>
	for _, strOut = range lstBattle {
		if strings.Contains(strOut, `<span>до начала `) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		// Проверка: если мы ещё в очереди, бой не начался
		for _, strOut = range lstBattle {
			if strings.Contains(strOut, `в очереди`) {
				fmt.Println("мастера: ещё в очереди, бой не начался")
				return
			}
		}
		// Не в очереди и нет таймера — бой идёт
		fmt.Println("мастера: бой идёт, переходим в режим ожидания")
		битва.Состояние().Уст(cons.РежимОжидание)
		return
	}
	// Найдена строка ожидания — выставляем таймер
	lstTime := strings.Split(strOut, `<span>до начала `)
	strTime := lstTime[1]
	lstTime = strings.Split(strTime, ` (`)
	strTime = lstTime[0]
	битва.ОбратВремяУст(АВремя(strTime))
}
