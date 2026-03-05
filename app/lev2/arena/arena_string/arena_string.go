// package arena_string -- потокобезопасный компонент списка строк для анализа арены
package arena_string

import (
	"strings"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// АренаСтроки -- потокобезопасный список строк арены
type АренаСтроки struct {
	val        []string
	strControl string // Контрольная строка в исходной строке для анализа арены
	block      sync.RWMutex
}

// НовАренаСтроки -- возвращает новый потокобезопасный список строк арены
func НовАренаСтроки(конт ILocalCtx, strControl string) *АренаСтроки {
	лог := NewLogBuf()
	лог.Debug("НовАренаСтроки(): strControl=%q", strControl)
	Hassert(strControl != "", "НовАренаСтроки(): strControl is empty")
	сам := &АренаСтроки{
		val:        make([]string, 0),
		strControl: strControl,
	}
	return сам
}

// Получ -- возвращает список строк для анализа
func (сам *АренаСтроки) Получ() []string {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// Set -- устанавливает список строк для анализа
func (сам *АренаСтроки) Set(lstString []string) {
	сам.block.Lock()
	defer сам.block.Unlock()
	if lstString == nil || len(lstString) == 0 {
		лог := NewLogBuf()
		лог.Warn("АренаСтроки.Set(): пустой список строк, ожидался title(%q) — пропускаем\n", сам.strControl)
		return
	}
	isOk := false
	for _, strControl := range lstString {
		if strings.Contains(strControl, сам.strControl) {
			isOk = true
			break
		}
	}
	if isOk {
		сам.val = lstString
		return
	}
	// Найти заголовок для отладки
	var strOut string
	for _, strOut = range lstString {
		if strings.Contains(strOut, "<title>") {
			break
		}
	}
	// Страница временно недоступна (например, идёт апгрейд здания) — пропускаем обновление
	лог := NewLogBuf()
	лог.Warn("АренаСтроки.Set(): ожидался title(%q), получен(%q) — пропускаем обновление\n", сам.strControl, strOut)
}
