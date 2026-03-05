// package arena_battle -- объект сражения
package arena_battle

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_battle/bf_battle_make"
	"wartank/app/lev2/arena/arena_battle/bf_battle_register"
	"wartank/app/lev2/arena/arena_battle/bf_battle_wait"
	"wartank/app/lev2/arena/arena_build"
)

// АренаСражение -- объект сражения
type АренаСражение struct {
	ИАренаСтроение
	конт   ILocalCtx
	клиент ИХттпВоркер
}

// НовСражение -- возвращает новую арену сражения PVE
func НовСражение(конт ILocalCtx) *АренаСражение {
	бот := конт.Get("бот").Val().(ИБот)
	сам := &АренаСражение{
		конт:   конт,
		клиент: бот.Сеть().ВебВоркер(),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Арена сражения",
		СтрКонтроль_: "<title>Сражения</title>",
		СтрУрл_:      "https://wartank.ru/pve",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)

	конт.Set("pve", сам, "Сражение с ботами")
	return сам
}

func (сам *АренаСражение) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_battle_register.СражениеРегистрация(сам.конт)
	bf_battle_wait.СражениеОжидать(сам.конт)
	bf_battle_make.СражениеВыполнить(сам.конт)
}
