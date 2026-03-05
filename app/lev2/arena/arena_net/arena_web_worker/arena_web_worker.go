// package arena_web_worker -- веб-воркер  арены
package arena_web_worker

import (
	"sync"

	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Объект сетевого соединения
*/

// АренаВебВоркер -- объект сетевого соединения
type АренаВебВоркер struct {
	ботСеть   ИБотСеть
	вебВоркер ИХттпВоркер
	лог       ILogBuf
	block     sync.Mutex
}

// НовАренаВебВоркер -- возвращает сетевого клиента
func НовАренаВебВоркер(конт ILocalCtx) *АренаВебВоркер { //, ботСеть ИБотСеть)
	лог := NewLogBuf()
	лог.Info("НовАренаВебВоркер()\n")
	сам := &АренаВебВоркер{
		ботСеть:   конт.Get("ботСеть").(ИБотСеть),
		вебВоркер: конт.Get("хттпВоркер").(ИХттпВоркер), //ботСеть.ВебВоркер(),
		лог:       лог,
	}
	return сам
}

// Получ -- выполняет безопасный GET-запрос в сеть
func (сам *АренаВебВоркер) Получ(ссылка string) (lstString []string, err error) {
	сам.block.Lock()
	defer сам.block.Unlock()
	// if ссылка == "https://wartank.ru/production/Mine" {
	сам.лог.Debug("АренаВебВоркер.Получ(): link=%v\n", ссылка)
	// }
	списОтвет := сам.вебВоркер.Получ(ссылка)
	return списОтвет, nil
}
