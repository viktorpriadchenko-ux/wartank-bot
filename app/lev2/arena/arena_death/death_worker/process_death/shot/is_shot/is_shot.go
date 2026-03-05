package is_shot

import (
	"sync"
)

/*
	Потокобезопасный признак внешнего разрешения выстрела
*/

// IsShot -- потокобезопасный признак внешнего разрешения выстрела
type IsShot struct {
	val   bool
	block sync.RWMutex
}

// NewIsShot -- возвращает новый *IsShot
func NewIsShot() *IsShot {
	return &IsShot{}
}

// Get -- возвращает хранимое состояние
func (сам *IsShot) Get() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// Set -- устанавливает хранимое состояние
func (сам *IsShot) Set() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = true
}

// Reset -- сбрасывает хранимое состояние
func (сам *IsShot) Reset() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = false
}
