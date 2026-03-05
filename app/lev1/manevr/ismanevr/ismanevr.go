package ismanevr

import (
	"sync"
)

/*
	Потокобезопасный признак внешнего разрешения манёвра
*/

// IsManevr -- потокобезопасный признак внешнего разрешения манёвра
type IsManevr struct {
	val   bool
	block sync.RWMutex
}

// NewIsManevr -- возвращает новый *IsManevr
func NewIsManevr() *IsManevr {
	return &IsManevr{
		val: true,
	}
}

// Get -- возвращает хранимое состояние
func (сам *IsManevr) Get() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// Set -- устанавливает хранимое состояние
func (сам *IsManevr) Set() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = true
}

// Сброс -- сбрасывает хранимое состояние
func (сам *IsManevr) Сброс() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = false
}
