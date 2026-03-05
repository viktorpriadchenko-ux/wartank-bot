package healthtime

import (
	"sync"
)

/*
	Содержит временное здоровье
*/

// HealthTime -- временное здоровье
type HealthTime struct {
	val   int // Здоровье FIXME: переделать в алиас
	block sync.RWMutex
}

// NewHealthTime -- возвращает новый *HealthTime
func NewHealthTime() *HealthTime {
	return &HealthTime{}
}

// Get -- возвращает хранимое временное здоровье
func (сам *HealthTime) Get() int {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// IsZero -- возвращает истину, если значение обнулено
func (сам *HealthTime) IsZero() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val == 0
}

// Set -- устанавливает значение по требованию
func (сам *HealthTime) Set(val int) {
	сам.block.Lock()
	defer сам.block.Unlock()
	if val < 0 {
		// log._rintf("WARN HealthTime.Set(): отрицательное значение(%v)\n", val)
		сам.val = 0
		return
	}
	сам.val = val
}
