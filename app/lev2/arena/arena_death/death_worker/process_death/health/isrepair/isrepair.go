package isrepair

import (
	"sync"
)

/*
	Потокобезопасный признак необходимости восстановления здоровья
*/

// IsRepair -- потокобезопасный признак восстановления здоровья
type IsRepair struct {
	val   bool
	block sync.RWMutex
}

// NewIsRepair -- возвращает новый *IsRepair
func NewIsRepair() *IsRepair {
	return &IsRepair{}
}

// Get -- возвращает хранимое состояние
func (сам *IsRepair) Get() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// Set -- устанавливает хранимое состояние
func (сам *IsRepair) Set() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = true
}

// Reset -- сбрасывает хранимое состояние
func (сам *IsRepair) Reset() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = false
}
