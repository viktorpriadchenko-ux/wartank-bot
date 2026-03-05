package parsehour

import (
	"fmt"
	"strconv"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
)

/*
	Потокобезопасно парсит из строки значение часа
*/

// ParseHour -- потокобезопасный парсер часа из строки
type ParseHour struct {
	val   int // Значение часа
	block sync.RWMutex
}

// NewParseHour -- возвращает новый *ParseHour
func NewParseHour() *ParseHour {
	return &ParseHour{}
}

// Get -- возвращает хранимое значение
func (sf *ParseHour) Get() int {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.val
}

// String -- возвращает строковое значение часов
func (sf *ParseHour) String() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	res := fmt.Sprintf("%02d", sf.val)
	return res
}

// Reset -- сбрасывает значение часов
func (sf *ParseHour) Reset() {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.val = 0
}

// Parse -- устанавливает значение часов
func (sf *ParseHour) Parse(strHour string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	iHour, err := strconv.Atoi(strHour)
	Hassert(err == nil, "ParseHour.Parse(): часы(%q) не число, err=\n\t%w", strHour, err)
	sf.set(iHour)
}

// Set - -устанавливает числовое значение часов
func (sf *ParseHour) Set(iHour int) {
	sf.block.Lock()
	defer sf.block.Unlock()
	sf.set(iHour)
}

// Внутренняя процедура для числовой установки часов без блокировки
func (sf *ParseHour) set(iHour int) {
	Hassert(iHour >= 0, "ParseHour.set(): часы(%v) < 0", iHour)
	sf.val = iHour
}
