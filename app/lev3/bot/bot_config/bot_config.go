// package bot_config -- конфиг бота для хранения в базе
package bot_config

import (
	"encoding/json"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"

	. "wartank/app/lev0/alias"
)

// БотКонфиг -- конфиг бота для хранения в базе
type БотКонфиг struct {
	ЕслиАвтозапуск_ bool      `json:"is_auto_run,omitempty"` // Признак автостарта при загрузке
	Логин_          string    `json:"login"`                 // Логин бота
	Пароль_         string    `json:"password"`              // Пароль бота
	Номер_          АБотНомер `json:"number"`                // Номер бота
	блок            sync.RWMutex
}

// Marshall -- сериализует конфиг в JSON
func (сам *БотКонфиг) Marshall() []byte {
	binData, _ := json.Marshal(сам)
	return binData
}

// Unmarshal -- десериализует себя из байтового потока
func (сам *БотКонфиг) Unmarshal(binData []byte) {
	лог := NewLogBuf()
	лог.Debug("БотКонфиг.Unmarshal()")
	err := json.Unmarshal(binData, сам)
	Hassert(err == nil, "Unmarshal(): err=\n\t%v\n", err)
}

// Логин -- возвращает логин
func (сам *БотКонфиг) Логин() string {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return сам.Логин_
}

// Номер -- возвращает номер бота
func (сам *БотКонфиг) Номер() АБотНомер {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return сам.Номер_
}
