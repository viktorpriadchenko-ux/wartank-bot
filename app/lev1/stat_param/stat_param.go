// package stat_param -- отдельный параметр статистики
package stat_param

import (
	"fmt"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// СтатПарам -- потокобезопасная статистика
type статПарам struct {
	sync.RWMutex
	имя  string // Имя параметра
	знач int    // Значение параметра
	лог  ILogBuf
}

// НовСтатПарам1 -- возвращает новую статистику
func НовСтатПарам1(имя string) ИСтатПарам {
	Hassert(имя != "", "НовСтатПарам(): имя параметра пустое\n")
	сам := &статПарам{
		имя: имя,
		лог: NewLogBuf(),
	}
	сам.лог.Info("НовСтатПарам(%s)\n", имя)
	return сам
}

// ЗначСтр -- строковое представление значение параметра
func (сам *статПарам) ЗначСтр() string {
	сам.RLock()
	defer сам.RUnlock()
	return fmt.Sprint(сам.знач)
}

// Уст -- устанавливает значение параметра
func (сам *статПарам) Уст(val int) {
	сам.Lock()
	defer сам.Unlock()
	сам.знач = val
}

// Получ -- значение хранимого параметра
func (сам *статПарам) Получ() int {
	сам.RLock()
	defer сам.RUnlock()
	return сам.знач
}

// Имя -- возвращает имя хранимого параметра
func (сам *статПарам) Имя() string {
	сам.RLock()
	defer сам.RUnlock()
	return сам.имя
}

// ИмяУст -- устанавливает значение имени
func (сам *статПарам) ИмяУст(val string) {
	сам.Lock()
	defer сам.Unlock()
	сам.имя = val
}
