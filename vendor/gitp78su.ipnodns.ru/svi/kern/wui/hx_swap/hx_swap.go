// package hx_swap -- атрибут HTMX (политика замены)
package hx_swap

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxSwap -- атрибут HTMX (политика замены)
type HxSwap struct {
	sync.RWMutex
	val string
}

// NewHxSwap -- возвращает новую политику замены
func NewHxSwap() *HxSwap {
	sf := &HxSwap{}
	_ = IHxSwap(sf)
	return sf
}

// String -- возвращает строковое представление тэга
func (sf *HxSwap) String() string {
	sf.RLock()
	defer sf.RUnlock()
	return `hx-swap="` + sf.val + `"`
}

// Get -- возвращает хранимое значение политики замена
func (sf *HxSwap) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает значение политики обмена
func (sf *HxSwap) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}
