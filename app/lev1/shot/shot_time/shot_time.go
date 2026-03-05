package shot_time

import (
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"

	. "wartank/app/lev0/alias"
)

/*
	Содержит время до выстрела
*/

// ShotTime -- время до выстрела
type ShotTime struct {
	val   АМилСек // Время в мсек
	block sync.RWMutex
}

// NewShotTime -- возвращает новый ShotTime
func NewShotTime() *ShotTime {
	return &ShotTime{}
}

// Get -- возвращает хранимое время до выстрела
func (сам *ShotTime) Get() АМилСек {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val
}

// Dec5 -- уменьшает время до выстрела на 5 мсек
func (сам *ShotTime) Dec5() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val -= 5
	if сам.val > 500_000 {
		сам.val = 0
	}
}

// Dec30 -- уменьшает время до выстрела на 30 мсек
func (сам *ShotTime) Dec30() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val -= 30
	if сам.val > 500_000 {
		сам.val = 0
	}
}

// Inc210 -- увеличивает время до выстрела на 210 мсек
func (сам *ShotTime) Inc210() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.val += 210
}

// IsZero -- возвращает истину, если значение обнулено
func (сам *ShotTime) IsZero() bool {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.val == 0
}

// Set -- устанавливает значение по требованию
func (сам *ShotTime) Set(val АМилСек) {
	сам.block.Lock()
	defer сам.block.Unlock()
	Hassert(val <= 500_000, "ShotTime.Set(): значение(%v)>500_000", val)
	сам.val = val
}
