package div_war_on

import (
	"context"
	"fmt"
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

// DivWarOn -- непосредственно танкует в сражении
type DivWarOn struct {
	ИАренаСтроение
	конт           ILocalCtx
	ctxDivWar      context.Context // Контекст сражения
	fnCancelDivWar func()          // Функция отмены сражения

	выстрел   ИВыстрел  // Объект выстрела
	health    ИЗдоровье // Текущее здоровье танка
	manevr    ИМанёвр   // Возможность маневрирования
	login     string
	isMasking ISafeBool // Признак необходимости маскирования (запрет стрельбы, когда слабое здоровье)
	chTick    chan int  // Ежесекундная проверка на окончание сражения
	isEnd     ISafeBool
}

// NewDivWarOn -- возвращает новый *DivWarOn
func NewDivWarOn(конт ILocalCtx) (*DivWarOn, error) {
	bot := конт.Get("бот").Val().(ИБот)
	ctxDivWar, fnCancelDivWar := context.WithTimeout(конт.Ctx(), time.Second*305)
	сам := &DivWarOn{
		конт:           конт,
		ctxDivWar:      ctxDivWar,
		fnCancelDivWar: fnCancelDivWar,
		login:          bot.Имя(),
		isMasking:      NewSafeBool(),
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
		return nil, fmt.Errorf("NewDivWarOn(): нет страницы для сражения")
	}
	go сам.run()
	return сам, nil
}

// Ежесекудный тик
func (сам *DivWarOn) makeTick() {
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

// Выстрел -- возвращает объект выстрела
func (сам *DivWarOn) Выстрел() ИВыстрел {
	return сам.выстрел
}

// запускает сражение
func (сам *DivWarOn) run() {
	// defer log._rintf("DivWarOn.run(): сражение завершено\n")
	{ // Подготовка к сражению
		сам.выстрел = shot.НовВыстрел(сам) // Объект выстрела
		сам.health = health.НовЗдоровье(сам)
		сам.manevr = manevr.НовМанёвр(сам)
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
func (сам *DivWarOn) checkEnd() bool {
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
func (сам *DivWarOn) ЕслиКонец() ISafeBool {
	return сам.isEnd
}

func (сам *DivWarOn) Манёвр() ИМанёвр {
	return сам.manevr
}

// ВыстрелБлок -- признак запрета стрельбы при слабом здоровье
func (сам *DivWarOn) ВыстрелБлок() ISafeBool {
	return сам.isMasking
}

// Ctx -- возвращает контекст отмены сражения
func (сам *DivWarOn) Ctx() context.Context {
	return сам.ctxDivWar
}

// CancelBattle - -вызов функции отмены контекста сражения
func (сам *DivWarOn) CancelBattle() {
	сам.fnCancelDivWar()
}
