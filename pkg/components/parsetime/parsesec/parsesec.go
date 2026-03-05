package parsesec

import (
	"fmt"
	"strconv"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
)

/*
	Парсер секунд
*/

// ParseSec -- парсер секунд
type ParseSec struct {
	val   int
	block sync.RWMutex
}

// NewParseSec -- возвращает новый *ParseSec
func NewParseSec() *ParseSec {
	return &ParseSec{}
}

// Get -- возвращает хранимое значение
func (sf *ParseSec) Get() int {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// String -- возвращает строковое значение секунд
func (sf *ParseSec) String() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	res := fmt.Sprintf("%02d", sf.val)
	return res
}

// Reset -- сбрасывает значение секунд
func (sf *ParseSec) Reset() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = 0
}

// Parse -- устанавливает значение секунд
func (sf *ParseSec) Parse(strSec string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	iSec, err := strconv.Atoi(strSec)
	Hassert(err == nil, "ParseSec.Parse(): секунды(%v) не число, err=%w", strSec, err)
	sf.set(iSec)
}

// Set -- устанавливает целочисленное значение
func (sf *ParseSec) Set(iSec int) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.set(iSec)
}

// Внутренняя установка значения секунд
func (sf *ParseSec) set(iSec int) {
	Hassert(0 <= iSec && iSec < 60, "ParseSec.set(): секунды(%v) не в диапазоне 0..60", iSec)
	sf.val = iSec
}
