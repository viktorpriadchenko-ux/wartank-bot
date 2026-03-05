// package server_stat -- глобальная статистика сервера
package server_stat

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	// . "wartank/app/lev0/types"
)

// СерверСтат -- структура статистики сервера
type СерверСтат struct {
	конт         IKernelCtx
	CчётСтарт_   int           `json:"count_start"`  // Количество запусков
	ВремяВсего_  time.Duration `json:"time_total"`   // Общее время работы в секундах
	ВремяСессия_ time.Duration `json:"time_session"` // Время сессии в секундах
	блок         sync.RWMutex
	лог          ILogBuf
}

// НовСерверСтат -- возвращает структуру статистики сервера
func НовСерверСтат() *СерверСтат {
	конт := GetKernelCtx()
	лог := NewLogBuf()
	лог.Info("НовСерверСтат()\n")
	сам := &СерверСтат{
		конт:         конт,
		CчётСтарт_:   0,
		ВремяВсего_:  0,
		ВремяСессия_: 0,
		лог:          лог,
	}
	сам.загр()
	сам.блок.Lock()
	сам.CчётСтарт_++
	сам.блок.Unlock()
	сам.сохр()
	go сам.пуск()
	return сам
}

// Загружает статистику сервера
func (сам *СерверСтат) загр() {
	store := сам.конт.Get("kernStoreKV").Val().(IKernelStoreKv)
	res := store.Get("server_stat")
	if res.IsErr() {
		if strings.Contains(res.Error().Error(), "not found") {
			return
		}
		Hassert(false, "СерверСтат.загр(): при загрузке статистики из хранилища, ош=\n\t%v\n", res.Error())
	}
	ош := json.Unmarshal(res.Unwrap(), сам)
	Hassert(ош == nil, "СерверСтат.загр(): при декодировании статистики из JSON, ош=\n\t%v\n", ош)
	go сам.пуск()
}

// Работает в отдельном потоке, считает время работы и счетчик запусков
func (сам *СерверСтат) пуск() {
	фнПуск := func() {
		сам.блок.Lock()
		defer сам.блок.Unlock()
		сам.ВремяВсего_ += time.Second
		сам.ВремяСессия_ += time.Second
		сам.сохр()
		time.Sleep(time.Second)
	}
	for {
		select {
		case <-сам.конт.Ctx().Done():
			return
		default:
			фнПуск()
		}
	}
}

// Сохраняет статистику сервера
func (сам *СерверСтат) сохр() {
	фнСохр := func() {
		сам.блок.Lock()
		defer сам.блок.Unlock()
		бинДанные, ош := json.Marshal(сам)
		Hassert(ош == nil, "СерверСтат.сохр(): при кодировании статистики в JSON, ош=\n\t%v\n", ош)
		store := сам.конт.Get("kernStoreKV").Val().(IKernelStoreKv)
		res := store.Set("server_stat", бинДанные)
		res.Hassert("СерверСтат.сохр(): при сохранении статистики в хранилище")
	}
	фнСохр()
}

// СчётСтарт -- счётчик запусков
func (сам *СерверСтат) СчётСтарт() int {
	return сам.CчётСтарт_
}

// ВремяСессия -- возвращает время сессии
func (сам *СерверСтат) ВремяСессия() string {
	return сам.ВремяСессия_.String()
}

// ВремяВсего -- возвращает общее время работы
func (сам *СерверСтат) ВремяВсего() string {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	return сам.ВремяВсего_.String()
}
