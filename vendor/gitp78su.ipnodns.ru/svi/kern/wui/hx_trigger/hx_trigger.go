// package hx_trigger -- атрибут HTMX (триггер запроса)
package hx_trigger

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxTrigger -- атрибут HTMX (триггер запроса)
type HxTrigger struct {
	sync.RWMutex
	val string
}

// NewHxTrigger -- возвращает новый триггер запроса
func NewHxTrigger() *HxTrigger {
	sf := &HxTrigger{}
	_ = IHxSwap(sf)
	return sf
}

// String -- возвращает строковое представление тэга
func (sf *HxTrigger) String() string {
	sf.RLock()
	defer sf.RUnlock()
	return `hx-trigger="` + sf.val + `"`
}

// Get -- возвращает хранимое значение триггера запроса
func (sf *HxTrigger) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает значение триггера запроса
func (sf *HxTrigger) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}
