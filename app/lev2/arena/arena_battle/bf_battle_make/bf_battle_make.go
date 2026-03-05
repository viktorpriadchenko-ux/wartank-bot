// package battle_worker -- бизнес-функция исполнения сражения
package bf_battle_make

import (
	"fmt"
	"strings"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena/arena_battle/bf_battle_make/battle_worker"
	"wartank/app/lev2/arena/arena_battle/bf_battle_make/battle_worker/battle_sound"
)

// НовСражениеДействие -- возвращает новый исполнитель битвы
func СражениеВыполнить(конт ILocalCtx) {
	сражение := конт.Get("pve").Val().(ИАренаСтроение)
	if сражение.Состояние().Получ() != cons.РежимОжидание {
		return
	}
	// Обновляем страницу и проверяем что это реальный бой
	сражение.Обновить()
	lstBattle := сражение.СписПолучить()
	еслиРеальныйБой := false
	for _, strOut := range lstBattle {
		// Ищем признак активного боя: кнопка атаки или HP-полоска
		if strings.Contains(strOut, `-currentControl-attack`) ||
			strings.Contains(strOut, `<span>Ремкомплект</span>`) ||
			strings.Contains(strOut, `-currentControl-maneuverLink`) {
			еслиРеальныйБой = true
			break
		}
	}
	if !еслиРеальныйБой {
		fmt.Printf("сражение: страница не содержит бой (строк=%d), сброс в построено\n", len(lstBattle))
		сражение.Состояние().Уст(cons.РежимПостроено)
		return
	}
	fmt.Println("сражение: подтверждён бой, запускаем горутины")
	пуск(конт)
}

// выполняет битву
func пуск(конт ILocalCtx) {
	действие := battle_worker.НовСражениеДействие(конт) // IBattleOn (онлайн)
	звук := battle_sound.NewBattleSound()
	звук.Play()
	time.Sleep(time.Second * 10) // Задержка для звука на странице
	<-действие.Контекст().Done()

	// Статистика: записать сражение (PVE)
	бот := конт.Get("бот").Val().(ИБот)
	бот.Статистика().СраженияДобавить(true) // TODO: определять исход по странице результата
}

// Тревога -- возвращает признак начала сражения (для браузера)
// func Тревога() ИСтатПарам {
// 	return сам.еслиНачало
// }
