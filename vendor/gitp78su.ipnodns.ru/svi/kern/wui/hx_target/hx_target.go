// package hx_target -- атрибут HTMX (цель замены)
package hx_target

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxSwap -- атрибут HTMX (цель замены)
type HxSwap struct {
	sync.RWMutex
	val string
}

// NewHxTarget -- возвращает новую цель замены
func NewHxTarget() *HxSwap {
	sf := &HxSwap{}
	_ = IHxTarget(sf)
	return sf
}

// String -- возвращает строковое представление тэга
func (sf *HxSwap) String() string {
	sf.RLock()
	defer sf.RUnlock()
	return `hx-target="` + sf.val + `"`
}

// Get -- возвращает хранимое значение цели замена
func (sf *HxSwap) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает значение цели обмена
func (sf *HxSwap) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}
