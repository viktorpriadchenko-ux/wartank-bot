// package arena_medal -- арена получения медалей
package arena_medal

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_medal/bf_medal_find"
)

// Банк -- объект арены медалей в миссиях
type АренаМедаль struct {
	ИАрена
	конт ILocalCtx
}

// НовБанк -- возвращает новую арену медалей в миссиях
func НовАренаМедали(конт ILocalCtx) *АренаМедаль {

	сам := &АренаМедаль{
		конт: конт,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Медали",
		СтрКонтроль_: `<title>Медали</title>`,
		СтрУрл_:      "https://wartank.ru/medals/current",
	}
	сам.ИАрена = arena.НовАрена(конт, аренаКонфиг)
	конт.Set("арена_медали", сам, "Арена получения медалей")
	return сам
}

func (сам *АренаМедаль) Пуск() {
	сам.ИАрена.Пуск()
	if сам.Состояние().Получ() == cons.РежимНеСуществует {
		сам.Состояние().Уст(cons.РежимПостроено)
	}
	bf_medal_find.МедалиНайти(сам.конт)
}
