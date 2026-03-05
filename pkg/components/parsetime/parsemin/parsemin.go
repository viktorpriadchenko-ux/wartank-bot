package parsemin

import (
	"fmt"
	"strconv"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
)

/*
	Потокобезопасно парсит из строки значение vbyen
*/

// ParseMin -- потокобезопасный парсер vbyen из строки
type ParseMin struct {
	val   int // Значение vbyen
	block sync.RWMutex
}

// NewParseMin -- возвращает новый *ParseMin
func NewParseMin() *ParseMin {
	return &ParseMin{}
}

// Get -- возвращает хранимое значение
func (sf *ParseMin) Get() int {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// String -- возвращает строковое значение минут
func (sf *ParseMin) String() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	res := fmt.Sprintf("%02d", sf.val)
	return res
}

// Reset -- сбрасывает значение минут
func (sf *ParseMin) Reset() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = 0
}

// Parse -- устанавливает значение минут
func (sf *ParseMin) Parse(strMin string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	iMin, err := strconv.Atoi(strMin)
	Hassert(err == nil, "ParseMin.Parse(): минуты(%v) не число, err=%w", strMin, err)
	sf.set(iMin)
}

// Set -- устанавливает целочисленное значение минут
func (sf *ParseMin) Set(iMin int) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.set(iMin)
}

// Внтренняя установка минут
func (sf *ParseMin) set(iMin int) {
	Hassert(0 <= iMin && iMin < 60, "ParseMin.set(): минуты не в диапазоне(%v) 0..60", iMin)
	sf.val = iMin
}
