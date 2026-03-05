package manevr

import (
	"context"
	"fmt"

	// "log"
	"strings"
	"time"

	"wartank/app/lev1/repair_time"
	// "wartank/internal/components/sound"
	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Пытается маневрировать после выстрела
*/

// манёвр -- маневрирует после выстрела
type манёвр struct {
	ИСражениеПроцесс // FIXME:
	лог              ILogBuf
	isEnd            ISafeBool
	ctxEnd           context.Context
	еслиНадо         ISafeBool
	еслиГотов        ISafeBool               // Возможность выполнить манёвр
	manevrTime       *repair_time.RepairTime // Время до восстановления манёвра
	chTick           chan int                // Тики для поиска маневра
}

// НовМанёвр -- возвращает новый манёвр
func НовМанёвр(проц ИСражениеПроцесс) ИМанёвр {
	Hassert(проц != nil, "НовМанёвр(): ИСражениеПроцесс==nil")
	сам := &манёвр{
		ИСражениеПроцесс: проц,
		лог:              проц.Бот().КонтБот().Log(),
		ctxEnd:           проц.Контекст(),
		еслиНадо:         NewSafeBool(),
		isEnd:            проц.ЕслиКонец(),
		еслиГотов:        NewSafeBool(),
		manevrTime:       repair_time.NewRepairTime(),
		chTick:           make(chan int, 1),
	}
	_ = сам.manevrTime.Уст("0") // При запуске боя есть возможность маневрировать
	go сам.makeTick()
	go сам.run()
	return сам
}

// ЕслиГотов -- возвращает признак готовности к манёвру
func (сам *манёвр) ЕслиГотов() bool {
	return сам.еслиГотов.Get()
}

// УстНадо -- устанавливает требование манёвра
func (сам *манёвр) УстНадо() {
	сам.еслиНадо.Set()
}

// Генерирует тик для уменьшения времени ожидания восстановления возможности манёвра
func (сам *манёвр) makeTick() {
	defer func() {
		close(сам.chTick)
		// log._rintf("манёвр.makeTick(): сражение завершено\n")
	}()
	for {
		select {
		case <-сам.ctxEnd.Done():
			return
		default:
			if сам.manevrTime.Получ() == 0 {
				сам.chTick <- 1
			}

			сам.manevrTime.Dec()
			time.Sleep(time.Second * 1)
		}
	}
}

// Рабочий цикл маневра (~)
func (сам *манёвр) run() {
	fmt.Println("маневр: run() горутина стартовала")
	for range сам.chTick {
		сам.findManevrTime()
		// Авто-выполнение манёвра когда запрошен и готов
		if сам.еслиНадо.Get() && сам.manevrTime.IsReady() {
			сам.Выполнить()
			сам.еслиНадо.Reset()
		}
	}
}

// Ищет время для манёвра
func (сам *манёвр) findManevrTime() {
	var (
		еслиНайдено bool
		ind         int
		lstBattleOn = сам.СписПолучить()
		стрВых      string
	)
	for ind, стрВых = range lstBattleOn {
		// <a href="pve?4-88.ILinkListener-currentControl-maneuverLink" class="simple-but blue"><span><span>5 секунд</span></span></a>
		if strings.Contains(стрВых, `-currentControl-maneuverLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Или манёвр успел восстановиться, или конец сражения
		if strings.Contains(стрВых, `<span>Маневр</span>`) {
			_ = сам.manevrTime.Уст("0")
			time.Sleep(time.Second * 1)
			return
		}
		if сам.isEnd.Get() {
			time.Sleep(time.Second * 1)
			return
		}
		сам.лог.Warn("findManevrTime(): ошибка в поиске времени манёвра, стрВых=%v", стрВых)
		time.Sleep(time.Second * 1)
		return
	}
	{ // Найти время манёвра
		lstTime := strings.Split(стрВых, `ILinkListener-currentControl-maneuverLink" class="simple-but blue"><span><span>`)
		if len(lstTime) != 2 {
			сам.лог.Err("манёвр.findManevrTime(): нет двух полей во времени ожидания инд=%v\n\n%v\n%v\n%v",
				ind,
				lstBattleOn[ind-1],
				стрВых,
				lstBattleOn[ind+1],
			)
			сам.еслиГотов.Reset()
			time.Sleep(time.Second * 1)
			return
		}
		strTime := lstTime[1]
		lstTime = strings.Split(strTime, ` секунд</span></span></a>`)
		strTime = lstTime[0]
		if err := сам.manevrTime.Уст(strTime); err != nil {
			сам.лог.Err("манёвр.findManevrTime(): при обновлении времени ожидания манёвра, ош=\n\t%v", err)
			сам.еслиГотов.Reset()
			time.Sleep(time.Second * 1)
			return
		}
	}
	сам.еслиГотов.Set()
	сам.лог.Info("манёвр.findManevrTime(): до манёвра, время=%v", сам.manevrTime.Получ())
}

// Выполнить -- принудительный манёвр по требованию
func (сам *манёвр) Выполнить() {
	var (
		еслиНайдено = false
		lstBattleOn = сам.СписПолучить()
		strOut      = ""
	)
	for _, strOut = range lstBattleOn {
		if strings.Contains(strOut, `<span>Маневр</span>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	lstLink := strings.Split(strOut, `<a href="`)
	if len(lstLink) < 2 {
		return
	}
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `" class="simple-but blue"><span><span>Маневр</span></span></a>`)
	strLink = "https://wartank.ru/" + lstLink[0]
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		return
	}
	lstBattleOn = res.Unwrap()
	сам.СтрОбновить(lstBattleOn)
	fmt.Println("маневр: выполнен")
}

// IsReady -- возвращает готовность манёвра
func (сам *манёвр) IsReady() bool {
	return сам.manevrTime.IsReady()
}
