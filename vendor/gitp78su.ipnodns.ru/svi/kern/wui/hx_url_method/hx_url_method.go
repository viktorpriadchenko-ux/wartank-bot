// package hx_url_method -- атрибут HTMX (метод запроса)
package hx_url_method

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxSwap -- атрибут HTMX (метод запроса)
type HxSwap struct {
	sync.RWMutex
	val string
}

// NewHxUrlMethod -- возвращает новый метод запроса
func NewHxUrlMethod() *HxSwap {
	sf := &HxSwap{
		val: "hx-post",
	}
	_ = IHxUrlMethod(sf)
	return sf
}

// Get -- возвращает хранимое значение метода запроса
func (sf *HxSwap) Get() string {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает значение метода запроса
func (sf *HxSwap) Set(val string) {
	sf.Lock()
	defer sf.Unlock()
	sf.val = val
}
