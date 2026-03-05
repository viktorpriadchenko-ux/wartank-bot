// package safe_bool_react -- потокобезопасный булевый признак с реакцией на своё изменение
package safe_bool_react

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// safeBoolReact -- потокобезопасный булевый признак с реакцией на своё изменение
type safeBoolReact struct {
	sync.RWMutex
	dict map[string]func(bool) // Словарь обратных вызовов
	val  bool
}

// NewSafeBoolReact -- возвращает новый потокобезопасный булевый признак с реакцией на своё изменение
func NewSafeBoolReact() ISafeBoolReact {
	sf := &safeBoolReact{
		dict: map[string]func(bool){},
	}
	return sf
}

// Delete -- удаляет функцию обратного вызова из наблюдения
func (sf *safeBoolReact) Delete(key string) {
	sf.Lock()
	defer sf.Unlock()
	delete(sf.dict, key)
}

// Add -- добавляет функцию обратного вызова
func (sf *safeBoolReact) Add(key string, fn func(bool)) {
	sf.Lock()
	defer sf.Unlock()
	Hassert(key != "", "safeBoolReact.Add(): key is empty")
	_, isOk := sf.dict[key]
	Hassert(!isOk, "safeBoolReact.Add(): key already exists")
	sf.dict[key] = fn
}

// Get -- возвращает хранимый булевый признак
func (sf *safeBoolReact) Get() bool {
	sf.RLock()
	defer sf.RUnlock()
	return sf.val
}

// Set -- устанавливает булевый признак
func (sf *safeBoolReact) Set() {
	sf.Lock()
	defer sf.Unlock()
	sf.val = true
	for _, fn := range sf.dict {
		fn(true)
	}
}

// Reset -- сбрасывает булевый признак
func (sf *safeBoolReact) Reset() {
	sf.Lock()
	defer sf.Unlock()
	sf.val = false
	for _, fn := range sf.dict {
		fn(false)
	}
}
