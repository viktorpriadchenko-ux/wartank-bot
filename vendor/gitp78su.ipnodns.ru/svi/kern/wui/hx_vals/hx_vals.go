// package hx_vals -- атрибут HTMX (словарь значений)
package hx_vals

import (
	"encoding/json"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// HxVals -- атрибут HTMX (словарь значений)
type HxVals struct {
	sync.RWMutex
	dict map[string]any
}

// NewHxVals -- возвращает новый словарь значений
func NewHxVals() *HxVals {
	sf := &HxVals{
		dict: map[string]any{},
	}
	_ = IHxVals(sf)
	return sf
}

// Len -- возвращает размер словаря
func (sf *HxVals) Len() int {
	sf.RLock()
	defer sf.RUnlock()
	return len(sf.dict)
}

// Del -- удаляет ключ словаря
func (sf *HxVals) Del(key string) {
	sf.Lock()
	defer sf.Unlock()
	delete(sf.dict, key)
}

// Clear -- очищает словарь значений
func (sf *HxVals) Clear() {
	sf.Lock()
	defer sf.Unlock()
	sf.dict = map[string]any{}
}

// String -- возвращает строковое представление тэга
func (sf *HxVals) String() string {
	sf.RLock()
	defer sf.RUnlock()
	binJson, _ := json.Marshal(sf.dict)
	return `hx-vals='` + string(binJson) + `'`
}

// Get -- возвращает хранимое значение словарь значений
func (sf *HxVals) Get(key string) any {
	sf.RLock()
	defer sf.RUnlock()
	return sf.dict[key]
}

// Set -- устанавливает значение словарь значений
func (sf *HxVals) Set(key string, val any) {
	sf.Lock()
	defer sf.Unlock()
	Hassert(key != "", "HxVals.Set(): key is empty")
	sf.dict[key] = val
}
