// package arena_net -- сетевая арена
package arena_net

import (
	"fmt"
	"strings"
	"sync"

	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// АренаСеть -- базовый тип для сетевых арен
type АренаСеть struct {
	ИБотСеть
	клиент ИХттпВоркер
	арена  ИАрена
	стрУрл string
	блок   sync.Mutex
	лог    ILogBuf
}

// НовАренаСеть -- возвращает новый *АренаСеть
func НовАренаСеть_(конт ILocalCtx, арена ИАрена, стрУрл string) *АренаСеть {
	Hassert(арена != nil, "НовАренаСеть(): ИСценаСтр == nil")
	Hassert(стрУрл != "", "НовАренаСеть(): стрУрл пустой\n")
	лог := NewLogBuf()
	лог.Info("НовАренаСеть(): strUrl=%q\n", стрУрл)
	сам := &АренаСеть{
		ИБотСеть: арена.Бот().Сеть(),
		арена:    арена,
		стрУрл:   стрУрл,
		клиент:   арена.Бот().Сеть().ВебВоркер(),
		лог:      лог,
	}
	_ = ИАренаСеть(сам)
	return сам
}

// Обновить -- обновляет список строк
func (сам *АренаСеть) Обновить() {
	сам.блок.Lock()
	defer сам.блок.Unlock()
	// FIXME: попытка разобраться, что за фигня творится
	// time.Sleep(time.Millisecond * 500)
	сам.лог.Debug("Обновить(): бот=%s\tсцена=%v\n", сам.арена.Бот().Имя(), сам.арена.Имя())
	lstString := сам.клиент.Получ(сам.стрУрл)
	сам.арена.СтрОбновить(lstString)
}

// Get -- выполняет GET-запрос по указанному URL
func (сам *АренаСеть) Get(strLink string) Result[[]string] {
	сам.блок.Lock()
	defer сам.блок.Unlock()
	// log._rintf("INFO АренаСеть.Get(): link=%v\n", сам.strUrl)
	if !strings.Contains(strLink, сам.стрУрл) {
		err := fmt.Errorf("АренаСеть.Get(): strLink(%v) не содержит strUrl(%v)", strLink, сам.стрУрл)
		return NewErr[[]string](err)
	}
	lstString := сам.клиент.Получ(strLink)
	return NewOk(lstString)
}
