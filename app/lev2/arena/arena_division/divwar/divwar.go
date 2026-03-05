package div_war

import (
	"strings"
	"sync"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	div_war_on "wartank/app/lev2/arena/arena_division/divwar/divwaron"
)

/*
	Объект ожидания и ведения битвы дивизий
*/

// DivWar -- объект ожидания битвы дивизий
type DivWar struct {
	ИАренаСтроение
	конт  ILocalCtx
	alarm ИСтатПарам

	// Непосредственная битва
	дивОн ИСражениеПроцесс
	block sync.Mutex

	chDivWar chan int // Сигнал начала битвы дивизий
}

// NewDivWar -- возвращает новый *DivWar
func NewDivWar(конт ILocalCtx) *DivWar {
	бот := конт.Get("бот").Val().(ИБот)
	сам := &DivWar{
		конт:     бот.КонтБот(),
		alarm:    lev1.НовСтатПарам("тревога"),
		chDivWar: make(chan int, 1),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        бот.КонтБот(),
		АренаИмя_:    "Битва дивизий",
		СтрКонтроль_: `<span>до начала `,
		СтрУрл_:      "https://wartank.ru/bitva",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(сам.конт, аренаКонфиг)
	go сам.run()
	go сам.резервТик()
	return сам
}

func (сам *DivWar) резервТик() {
	for {
		select {
		case <-сам.конт.Ctx().Done():
			return
		default:
			ct0 := сам.ВремяОстат().ПолучМилСек()
			time.Sleep(time.Second * 7)
			ct1 := сам.ВремяОстат().ПолучМилСек()
			if ct1.Сек() != ct0.Сек() {
				continue
			}
			if сам.дивОн != nil {
				continue
			}
			сам.chDivWar <- 1
		}
	}
}

// запускает в работу битву дивизий
func (сам *DivWar) run() {
	сам.chDivWar <- 1
	for {
		select {
		case <-сам.конт.Ctx().Done():
			return
		case <-сам.ВремяОстат().КаналСиг(): // Время обновить данные по сражению
			сам.findTimeCount()
			сам.upDivWar()
		case <-сам.chDivWar: // Сигнал к началу сражения
			сам.block.Lock()
			if сам.дивОн != nil {
				continue
			}
			сам.alarm.Уст(1)
			go сам.DivWar() // Запустить цикл непосредственного сражения
			time.Sleep(time.Second * 5)
			сам.alarm.Уст(0)
		}
	}
}

// Ищет время до начала битвы дивизий
func (сам *DivWar) findTimeCount() {
	сам.Обновить()
	var (
		strOut      string
		lstDivWar   = сам.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstDivWar {
		if strings.Contains(strOut, `до начала: `) {
			ind++
			strOut = lstDivWar[ind]
			еслиНайдено = true
			break
		}
		if strings.Contains(strOut, `>ОБЫЧНЫЕ<`) { // Это уже битва
			сам.chDivWar <- 1
			return
		}
	}
	if !еслиНайдено { // Битва дивизий уже идёт
		сам.chDivWar <- 1
		return
	}
	lstTime := strings.Split(strOut, `<span>`)
	strTime := lstTime[1]
	lstTime = strings.Split(strTime, `</span>`)
	strTime = lstTime[0]
	сам.ОбратВремяУст(АВремя(strTime))
}

// При необходимости поднимает взвод в атаку, вызывается только если обнаружено приглашение (+)
func (сам *DivWar) upDivWar() {
	var (
		strOut      string
		lstDivWar   = сам.СписПолучить()
		еслиНайдено bool
	)
	for _, strOut = range lstDivWar {
		if strings.Contains(strOut, `>Взвод, подъем! В атаку!<`) {
			еслиНайдено = true
			break
		}
		if strings.Contains(strOut, `<div class="white medium cntr bold mb5">Вы в рядах участников</div>`) {
			// log._rintf("INFO DivWar.upDivWar(): уже зарегистрирован\n")
			return
		}
	}
	if !еслиНайдено {
		return
	}
	// Найдено приглашение на участие
	lstUp := strings.Split(strOut, `<a class="simple-but border" href="`)
	linkUp := lstUp[1]
	lstUp = strings.Split(linkUp, `"><span><span>Взвод, подъем! В атаку!</span></span></a>`)
	linkUp = "https://wartank.ru/" + lstUp[0]
	res := сам.Сеть().Get(linkUp)
	if res.IsErr() {
		return
	}
	сам.СтрОбновить(res.Unwrap())
}

// Ведёт сражение
func (сам *DivWar) DivWar() {
	defer func() {
		сам.дивОн = nil
		сам.block.Unlock()
		сам.ОбратВремяУст("01")
	}()
	дивОн, err := div_war_on.NewDivWarOn(сам.конт)
	if err != nil {
		return // Битвы ещё нет
	}
	сам.дивОн = дивОн
	// Цикл ожидания окончания сражения
	for !сам.дивОн.ЕслиКонец().Get() {
		time.Sleep(time.Second * 1)
	}
}

// Alarm -- возвращает признак начала сражения (для браузера)
func (сам *DivWar) Alarm() ИСтатПарам {
	return сам.alarm
}
