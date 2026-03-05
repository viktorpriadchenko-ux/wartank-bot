// package web_log -- веб-лог компонента
package web_log

import (
	"fmt"
	"strings"
	"sync"
	"time"
	. "wartank/app/lev0/types"
)

// ВебЛог -- веб-лог компонента
type ВебЛог struct {
	лог     []string // Буфер логирования
	еслиПеч bool
	блок    sync.RWMutex
}

// НовВебЛог -- возвращает новый *ВебЛог
func НовВебЛог(еслиПеч bool) ИВебЛог {
	сам := &ВебЛог{
		лог:     []string{},
		еслиПеч: еслиПеч,
	}
	_ = ИВебЛог(сам)
	return сам
}

// ОтклВывод -- отключает вывод в консоль
func (сам *ВебЛог) ОтклВывод() {
	сам.блок.Lock()
	defer сам.блок.Unlock()
	сам.еслиПеч = false
}

// Добавить -- добавляет строку в лог
func (сам *ВебЛог) Добавить(сбщ string, арг ...interface{}) {
	сам.блок.Lock()
	defer сам.блок.Unlock()
	время := time.Now().Local().Format("2006-01-02 15:04:05.000 ")
	сбщ = время + fmt.Sprintf(сбщ, арг...)
	if сам.еслиПеч {
		fmt.Print(сбщ)
	}
	сам.лог = append(сам.лог, сбщ)
	if len(сам.лог) > 50 {
		сам.лог = сам.лог[50:]
	}
}

// Ошибка -- лог ошибок
func (сам *ВебЛог) Ошибка() string {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return strings.Join(сам.лог, "\n")
}

// Внимание -- лог предупреждений
func (сам *ВебЛог) Внимание() string {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return strings.Join(сам.лог, "\n")
}

// Инфо -- информационный лог
func (сам *ВебЛог) Инфо() string {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return strings.Join(сам.лог, "\n")
}

// Отладка -- возвращает весь лог компонента
func (сам *ВебЛог) Отладка() string {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return strings.Join(сам.лог, "\n")
}
