package arena_missions

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_missions/bf_mission_simple"
)

/*
	Сканирует миссии на предмет забрать золотишко.
*/

// АренаМиссии -- забирает золотишко
type АренаМиссии struct {
	ИАрена
	конт ILocalCtx
	бот  ИБот
}

// НовМиссии -- возвращает новый *Миссии
func НовМиссии(конт ILocalCtx) *АренаМиссии {
	сам := &АренаМиссии{
		конт: конт,
		бот:  конт.Get("бот").Val().(ИБот),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Миссии",
		СтрКонтроль_: "<title>Миссии</title>",
		СтрУрл_:      "https://wartank.ru/missions/", // https://wartank.ru/missions/
	}
	сам.ИАрена = arena.НовАрена(конт, аренаКонфиг)
	конт.Set("миссии_простые", сам, "Арена простых миссий")
	_ = ИАренаМиссииПростые(сам)
	return сам
}

func (сам *АренаМиссии) Пуск() {
	сам.ИАрена.Пуск()
	bf_mission_simple.МиссииПростыеЗабрать(сам.конт)
}
