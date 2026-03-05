package is_play

import (
	"sync"
)

/*
	Потокобезопасный признак проигрывания звука
*/

// IsPlay -- потокобезопасный признак проигрывания звука
type IsPlay struct {
	val   bool
	block sync.RWMutex
}

// NewIsPlay -- возвращает новый *IsPlay
func NewIsPlay() *IsPlay {
	return &IsPlay{}
}

// Get -- возвращает хранимое состояние
func (сам *IsPlay) Get() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// Set -- устанавливает хранимое состояние
func (сам *IsPlay) Set() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = true
}

// Reset -- сбрасывает хранимое состояние
func (сам *IsPlay) Reset() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val = false
}
