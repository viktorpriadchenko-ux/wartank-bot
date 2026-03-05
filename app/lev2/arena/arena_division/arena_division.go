package arena_division

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev2/arena/arena_division/divwar"
)

/*
	Дивизия и все объекты в ней
*/

// АренаДивизия -- дивизия и все объекты в ней
type АренаДивизия struct {
	война *div_war.DivWar
}

// НовДивизия -- возвращает новую арену дивизии (запускает фоновые горутины)
func НовДивизия(конт ILocalCtx) *АренаДивизия {
	сам := &АренаДивизия{
		война: div_war.NewDivWar(конт),
	}
	return сам
}
