// package death_worker -- исполнение схватки
package death_worker

import (
	"time"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_death/death_worker/process_death"
	"wartank/app/lev2/arena/arena_death/death_worker/process_death/battle_sound"
)

// СражениеДействие -- исполнение схватки
type СхваткаИсполнитель struct {
	ИАрена
	конт ILocalCtx

	еслиНачало ИСтатПарам

	// Непосредственное сражение
	действие *process_death.СхваткаПроцесс

	sound *battle_sound.BattleSound // Однопоточное проигрывание звука
}

// НовСражениеДействие -- возвращает новый исполнитель схватки
func НовСхваткаИсполнитель(конт ILocalCtx) *СхваткаИсполнитель {
	сам := &СхваткаИсполнитель{
		конт:       конт,
		еслиНачало: lev1.НовСтатПарам("тревога"),
		sound:      battle_sound.NewBattleSound(),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Ход сражения",
		СтрКонтроль_: `<title>Сражения</title>`,
		СтрУрл_:      "https://wartank.ru/pve",
	}
	сам.ИАрена = arena.НовАрена(конт, аренаКонфиг)
	return сам
}

// Танковать -- выполняет битву
func (сам *СхваткаИсполнитель) Танковать() {
	сам.действие = process_death.НовСхваткаПроцесс(сам.конт) // IBattleOn (онлайн)
	сам.sound.Play()
	time.Sleep(time.Second * 10) // Задержка для звука на странице
	<-сам.действие.Контекст().Done()
	// log._rintf("Battle.runBaton(): сражение завершено\n")
}

// Тревога -- возвращает признак начала сражения (для браузера)
func (сам *СхваткаИсполнитель) Тревога() ИСтатПарам {
	return сам.еслиНачало
}
