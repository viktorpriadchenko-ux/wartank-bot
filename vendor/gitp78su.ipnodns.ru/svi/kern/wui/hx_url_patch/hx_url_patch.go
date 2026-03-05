// package hx_url_patch -- атрибут HTMX (путь запроса)
package hx_url_patch

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxUrlPatch -- атрибут HTMX (путь запроса)
type HxUrlPatch struct {
	sync.RWMutex
	val string
}

// NewHxUrlPatch -- возвращает новый путь запроса
func NewHxUrlPatch(patch string) *HxUrlPatch {
	Hassert(patch != "", "NewHxUrlPatch(): patch isempty")
	sf := &HxUrlPatch{
		val: patch,
	}
	_ = IHxUrlMethod(sf)
	return sf
}

// Get -- возвращает хранимое значение пути запроса
func (sf *HxUrlPatch) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает значение пути запроса
func (sf *HxUrlPatch) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}
