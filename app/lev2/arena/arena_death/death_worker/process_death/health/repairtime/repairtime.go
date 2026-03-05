package repairtime

import (
	"fmt"
	"strconv"
	"sync"
)

/*
	Потокобезопасное время ожидания до восстановлении ремки
*/

// RepairTime -- потокобезопасное время до восстановления ремки
type RepairTime struct {
	val    int
	valOld int // Старое время до ремки
	block  sync.RWMutex
}

// NewRepairTime -- возвращает новый *RepairTime
func NewRepairTime() *RepairTime {
	return &RepairTime{}
}

// Получ -- возвращает хранимое значение времени
func (сам *RepairTime) Получ() int {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// GetOld -- возвращает хранимое старое значение времени
func (сам *RepairTime) GetOld() int {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.valOld
}

// Set -- устанавливает хранимое время восстановления ремки
func (сам *RepairTime) Set(val string) error {
	сам.block.Lock()
	defer сам.block.Unlock()
	iVal, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("RepairTime.Set(): val(%v) не число, err=%w", val, err)
	}
	if iVal < 0 {
		return fmt.Errorf("RepairTime.Set(): val(%v) < 0", iVal)
	}
	сам.valOld = сам.val
	сам.val = iVal
	return nil
}

// Dec -- уменьшает на секунду время восстановления
func (сам *RepairTime) Dec() {
	сам.block.Lock()
	defer сам.block.Unlock()
	if сам.val > 0 {
		сам.valOld = сам.val
		сам.val--
	}
}

// ЕслиГотово -- возвращает признак готовности восстановления
func (сам *RepairTime) ЕслиМожно() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val == 0
}

// ЕслиИзменилось -- возвращает признак изменения здоровья после присвоения
func (сам *RepairTime) ЕслиИзменилось() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val == сам.valOld
}
