package damage

import (
	"sync"

	"wartank/app/lev0/alias"
)

/*
	Следит за уроном танка.
	Сравнивает с предыдущим значением,
	результат сравнения при каждом присвоении -- сохраняет
*/

// Damage -- урон танка с памятью
type Damage struct {
	val   alias.АУрон
	res   string
	block sync.RWMutex
}

// NewDamage -- возвращает новый *Damage
func NewDamage() *Damage {
	return &Damage{
		res: "none",
	}
}

// Set -- устанавливает урон
func (сам *Damage) Set(val alias.АУрон) {
	сам.block.Lock()
	defer сам.block.Unlock()
	if сам.val == 0 { // Первоначальное присвоение
		сам.val = val
		return
	}
	switch {
	case сам.val > val: // Урон уменьшился
		сам.res = "down"
	case сам.val < val: // Урон увеличился
		сам.res = "up"
	default: // Урон не изменился
		сам.res = "none"
	}
	сам.val = val
}

// Result -- возвращает результат сравнения урона старого и текущего
func (сам *Damage) Result() string {
	сам.block.RLock()
	defer сам.block.RUnlock()
	return сам.res
}
