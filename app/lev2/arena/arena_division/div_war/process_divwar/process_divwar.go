package process_divwar

import (
	"context"
	"strings"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/health"
	"wartank/app/lev1/manevr"
	"wartank/app/lev1/shot"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Предоставляет сетевой компонент при непосредственном сражении
*/

// ПроцессДивизияВойна -- непосредственно танкует в сражении
type ПроцессДивизияВойна struct {
	ИАренаСтроение
	конт           ILocalCtx
	лог            ILogBuf
	ctxDivWar      context.Context // Контекст сражения
	fnCancelDivWar func()          // Функция отмены сражения

	выстрел  ИВыстрел  // Объект выстрела
	здоровье ИЗдоровье // Текущее здоровье танка
	манёвр   ИМанёвр   // Возможность маневрирования
	логин    string
	chTick   chan int // Ежесекундная проверка на окончание сражения
	isEnd    ISafeBool
}

// НовПроцессДивизияВойна -- возвращает новый *DivWarOn
func НовПроцессДивизияВойна(конт ILocalCtx) ИСражениеПроцесс {
	лог := NewLogBuf()
	бот := конт.Get("бот").Val().(ИБот)
	ctxDivWar, fnCancelDivWar := context.WithTimeout(бот.КонтБот().Ctx(), time.Second*305)
	сам := &ПроцессДивизияВойна{
		лог:            лог,
		ctxDivWar:      ctxDivWar,
		fnCancelDivWar: fnCancelDivWar,
		логин:          бот.Имя(),
		isEnd:          NewSafeBool(),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Сражение",
		СтрКонтроль_: `<title>Сражения</title>`,
		СтрУрл_:      "https://wartank.ru/pve",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	if сам.checkEnd() {
		return nil
	}
	go сам.makeTick()
	go сам.run()
	return сам
}

// Манёвр -- возвращает объект манёвра
func (сам *ПроцессДивизияВойна) Манёвр() ИМанёвр {
	return сам.манёвр
}

// Выстрел -- возвращает объект выстрела
func (сам *ПроцессДивизияВойна) Выстрел() ИВыстрел {
	return сам.выстрел
}

// Ежесекудный тик
func (сам *ПроцессДивизияВойна) makeTick() {
	defer func() {
		close(сам.chTick)
		сам.isEnd.Set()
	}()
	for !сам.isEnd.Get() {
		select {
		case <-сам.конт.Ctx().Done(): // Отмена контекста приложения
			return
		case <-сам.ctxDivWar.Done(): // Битва закончилась
			return
		default:
			if сам.isEnd.Get() {
				return
			}
			сам.chTick <- 1
			time.Sleep(time.Second * 1)
		}
	}
}

// запускает сражение
func (сам *ПроцессДивизияВойна) run() {
	// defer log._rintf("ПроцессДивизияВойна.run(): сражение завершено\n")
	{ // Подготовка к сражению
		сам.выстрел = shot.НовВыстрел(сам) // Объект выстрела
		сам.здоровье = health.НовЗдоровье(сам)
		сам.манёвр = manevr.НовМанёвр(сам)
	}
	for { // Рабочий цикл сражения
		select {
		case <-сам.ctxDivWar.Done():
			return
		case <-сам.ВремяОстат().КаналСиг():
			if сам.checkEnd() {
				return
			}
		}
	}
}

// Проверяет окончание сражения
func (сам *ПроцессДивизияВойна) checkEnd() bool {
	defer func() {
		if сам.isEnd.Get() {
			сам.fnCancelDivWar()
			// log._rintf("DivWarOn.checkEnd(): сражение завершено\n")
		}
	}()

	сам.Обновить()
	lstDivWarOn := сам.СписПолучить()
	for _, strOut := range lstDivWarOn {
		if strings.Contains(strOut, `" class="simple-but gray"><span><span>ОБЫЧНЫЕ</span></span></a>`) {
			сам.isEnd.Reset()
			return false
		}
	}
	сам.isEnd.Set()
	сам.fnCancelDivWar()
	return true
}

// ЕслиКонец -- возвращает признак окончания сражения (интерфейс)
func (сам *ПроцессДивизияВойна) ЕслиКонец() ISafeBool {
	return сам.isEnd
}

// Ctx -- возвращает контекст отмены сражения
func (сам *ПроцессДивизияВойна) Ctx() context.Context {
	return сам.ctxDivWar
}

// CancelBattle - -вызов функции отмены контекста сражения
func (сам *ПроцессДивизияВойна) CancelBattle() {
	сам.fnCancelDivWar()
}
