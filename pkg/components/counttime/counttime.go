package counttime

import (
	"sync"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/pkg/components/parsetime"
)

/*
	Счётчик обратного временив мсек
*/

const (
	iSleep = time.Millisecond * 100
)

// CountTime -- счётчик обратного времени
type CountTime struct {
	kCtx   IKernelCtx
	val    ISafeInt
	parser *parsetime.ParseTime

	chTick     chan int
	chCall     chan int  // Канал для отправки сигналов
	isWork     ISafeBool // Признак работы
	timeTarget ISafeInt  // Целевое время срабатывания

	block sync.RWMutex
}

// NewCountTime -- возвращает новый *CountTime
func NewCountTime() *CountTime {
	sf := &CountTime{
		kCtx:       GetKernelCtx(),
		val:        NewSafeInt(),
		chTick:     make(chan int, 3),
		chCall:     make(chan int, 2),
		isWork:     NewSafeBool(),
		parser:     parsetime.NewParseTime(),
		timeTarget: NewSafeInt(),
	}
	sf.isWork.Set()
	go sf.makeTick()
	go sf.run()
	return sf
}

// Запускает тикер для секундных интервалов
func (sf *CountTime) makeTick() {
	defer func() {
		if !sf.isWork.Get() {
			return
		}
		sf.isWork.Reset()
		close(sf.chTick)
		// log._rintf("CountTime.makeTick(): работа завершена")
	}()
	countSleep := time.Duration(0)
	for {
		select {
		case <-sf.kCtx.Ctx().Done(): // Отмена контекста приложения
			// log._rintf("CountTime.makeTick(): глобальная отмена контекста\n")
			return
		default:
			if !sf.isWork.Get() {
				return
			}

			if countSleep >= time.Second {
				sf.chTick <- 1
				countSleep = 0
			}
			countSleep += iSleep
			time.Sleep(iSleep)
		}
	}
}

// Главный цикл обратного отсчёта
func (sf *CountTime) run() {
	for range sf.chTick {
		timeNow := time.Now().UTC().Unix()
		if sf.timeTarget.Get() > int(timeNow) {
			val := sf.timeTarget.Get() - int(timeNow)
			val -= 1
			sf.Set(val)
			continue
		}
		sf.chCall <- 1
		continue
	}
	close(sf.chCall)
}

// Stop -- останавливает работу счётчика
func (sf *CountTime) Stop() {
	sf.isWork.Reset()
}

// Get -- возвращает число оставшихся сек
func (sf *CountTime) Get() int {
	return sf.val.Get()
}

// Parse -- устанавливает число оставшихся сек
func (sf *CountTime) Parse(val string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(val != "", "CountTime.Set(): val is empty")
	sf.parser.Parse(val)
	_val := sf.parser.Hour().Get()*3600 + sf.parser.Min().Get()*60 + sf.parser.Min().Get()
	sf.val.Set(_val)
	_val = int(time.Now().UTC().Unix()) + sf.val.Get()
	sf.timeTarget.Set(_val)
}

// Set -- устанавливает число оставшихся сек
func (sf *CountTime) Set(val int) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(val >= 0, "CountTime.Set(): val(%v)<0", val)
	sf.val.Set(val)
	{ // Обновить локальные счётчики
		if val < 60 {
			sf.parser.Hour().Reset()
			sf.parser.Min().Reset()
			sf.parser.Sec().Set(val)
			return
		}
		if 60 < val && val < 3600 {
			sf.parser.Hour().Reset()
			iMin := val / 60
			sf.parser.Min().Set(iMin)
			val -= iMin * 60
			sf.parser.Sec().Set(val)
			return
		}
		sf.parser.Hour().Set(val / 3600)
		val -= sf.parser.Hour().Get() * 3600
		sf.parser.Min().Set(val / 60)
		val -= sf.parser.Min().Get() * 60
		sf.parser.Sec().Set(val)
		val = int(time.Now().UTC().Unix()) + sf.val.Get()
		sf.timeTarget.Set(val)
	}
}

// String -- возвращает строковое представление оставшихся сек
func (sf *CountTime) String() string {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.parser.String()
}

// ChanSig -- возвращает канал чтения тиков
func (sf *CountTime) ChanSig() <-chan int {
	return sf.chCall
}
