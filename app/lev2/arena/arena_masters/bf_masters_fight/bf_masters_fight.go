// package bf_masters_fight -- исполнение битвы мастеров
package bf_masters_fight

import (
	"fmt"
	"strings"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena/arena_masters/bf_masters_fight/masters_worker"
)

// МастераВыполнить -- запускает ведение битвы мастеров, если пора
func МастераВыполнить(конт ILocalCtx) {
	битва := конт.Get("pvp").Val().(ИАренаСтроение)
	/*
		Запускаем бой только когда:
		  - Таймер истёк и идёт сама битва (РежимОжидание)
	*/
	if битва.Состояние().Получ() != cons.РежимОжидание {
		return
	}
	// Обновляем страницу и проверяем что это реальный бой, а не лобби/очередь
	битва.Обновить()
	lstBattle := битва.СписПолучить()
	еслиРеальныйБой := false
	for _, strOut := range lstBattle {
		if strings.Contains(strOut, `-currentControl-attack`) ||
			strings.Contains(strOut, `<span>Ремкомплект</span>`) ||
			strings.Contains(strOut, `-currentControl-maneuverLink`) {
			еслиРеальныйБой = true
			break
		}
	}
	if !еслиРеальныйБой {
		fmt.Printf("мастера: страница не содержит бой (строк=%d), сброс в построено\n", len(lstBattle))
		битва.Состояние().Уст(cons.РежимПостроено)
		return
	}
	fmt.Println("мастера: подтверждён бой, запускаем горутины")
	пуск(конт)
	// После боя сбрасываем в "построено", чтобы в следующем цикле заново зарегистрироваться
	битва.Состояние().Уст(cons.РежимПостроено)
}

// пуск -- ведёт фактический бой (блокирует до завершения, макс ~5 минут)
func пуск(конт ILocalCtx) {
	действие := masters_worker.НовМастераДействие(конт)
	time.Sleep(time.Second * 10) // Задержка на подгрузку страницы
	<-действие.Контекст().Done()

	// Статистика: записать битву мастеров (PVP)
	бот := конт.Get("бот").Val().(ИБот)
	бот.Статистика().БитваМастеровДобавить(true) // TODO: определять исход
}
