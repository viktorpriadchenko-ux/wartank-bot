// package hx_swap_oob -- объект внеполосной подкачки
package hx_swap_oob

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxSwapOob -- объект внеполосной подкачки
type HxSwapOob struct {
	sync.RWMutex
	val string
}

// NewHxSwapOob -- возвращает новую внеполосную подкачку
func NewHxSwapOob() *HxSwapOob {
	sf := &HxSwapOob{}
	_ = IHxSwapOob(sf)
	return sf
}

// String -- возвращает строковое представление тэга
func (sf *HxSwapOob) String() string {
	sf.RLock()
	defer sf.RUnlock()
	return `hx-swap-oob="` + sf.val + `"`
}

// Get -- возвращает хранимое значение внеполосной подкачки
func (sf *HxSwapOob) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает значение внеполосной подкачки
func (sf *HxSwapOob) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}
