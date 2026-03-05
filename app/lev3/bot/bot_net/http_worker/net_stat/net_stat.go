// package net_stat -- сетевая статистика по обмену
package net_stat

import (
	"sync"
	"time"

	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Сетевая статистика для сетевого соединения
*/

// СетьСтата -- статистика сетевого соединения
type СетьСтата struct {
	ботСеть              ИБотСеть
	countByteInterval    int      // Счётчик байтов в текущем интервале
	totalByte            int      // Сколько всего байт передано
	всегоМинут           int      // Сколько всего минут работала передача
	countRequestInterval int      // Число запросов
	countErr             int      // Число зафиксированных ошибок
	chTick               chan int // Сигналы времени 1 раз в 5 минут
	block                sync.Mutex
	лог                  ILogBuf
}

// НовСетьСтата -- возвращает новый *NetStat
func НовСетьСтата(конт ILocalCtx) *СетьСтата {
	лог := NewLogBuf()
	лог.Info("НовСетьСтата()\n")

	ботСеть := конт.Get("ботСеть").Val().(ИБотСеть)
	сам := &СетьСтата{
		ботСеть: ботСеть,
		chTick:  make(chan int, 2),
		лог:     лог,
	}
	go сам.пуск()
	go сам.makeTick()
	return сам
}

// Тикер меток времени 1 раз в 5 минут
func (сам *СетьСтата) makeTick() {
	defer close(сам.chTick)
	for {
		select {
		case <-сам.ботСеть.Контекст().Ctx().Done():
			return
		default:
			time.Sleep(time.Second * 5 * 60)
			сам.chTick <- 1
		}
	}
}

// Главный цикл работы статистики
func (сам *СетьСтата) пуск() {
	fnCalc := func() {
		сам.block.Lock()
		defer сам.block.Unlock()
		сам.всегоМинут += 5
		_МБ := float32(сам.totalByte) / float32(сам.всегоМинут*60) * (3600 * 24 * 30.5) / (1024 * 1024)
		сам.лог.Info("пуск().fnCalc(): запросы=%0.2f/сек\t траф0=%0.2f бит/сек\tтраф1=%0.2f МБ/мес\tошибки=%v\n",
			float32(сам.countRequestInterval)/300,
			float32(сам.countByteInterval*8)/300,
			_МБ,
			сам.countErr)
		сам.countByteInterval = 0
		сам.countRequestInterval = 0
	}
	for range сам.chTick {
		select {
		case <-сам.ботСеть.Контекст().Ctx().Done():
			return
		case <-сам.chTick:
			fnCalc()
		}
	}
}

// IncErr -- добавляет ошибку в статистику
func (сам *СетьСтата) IncErr() {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.countErr++
}

// AddByte -- увеличивает счётчик запросов и байтов на передачу/приём
func (сам *СетьСтата) AddByte(val int) {
	сам.block.Lock()
	defer сам.block.Unlock()
	сам.countByteInterval += val
	сам.totalByte += val
	сам.countRequestInterval++
}
