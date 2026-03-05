package div_war

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_division/div_war/div_war_net"
	"wartank/app/lev2/arena/arena_division/div_war/process_divwar"
	"wartank/app/lev2/arena/arena_division/div_war/process_divwar/div_war_sound"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Объект ожидания битвы дивизий
*/

// DivWar -- объект ожидания битвы дивизий
type DivWar struct {
	ИАренаСтроение
	конт  ILocalCtx
	alarm ИСтатПарам
	net   *div_war_net.DivWarNet
	conn  ИХттпВоркер

	// Непосредственная битва
	процВойна ИСражениеПроцесс
	логин     string // Для непосредственной битвы дивизий
	block     sync.Mutex

	chDivWar chan int // Сигнал начала битвы дивизий

	sound *div_war_sound.DivWarSound // Однопоточное проигрывание звука
}

// NewDivWar -- возвращает новый *DivWar
func NewDivWar(конт IKernelCtx) (*DivWar, error) {
	bot := конт.Get("бот").Val().(ИБот)
	сам := &DivWar{
		конт:     конт,
		alarm:    lev1.НовСтатПарам("тревога"),
		chDivWar: make(chan int, 1),
		sound:    div_war_sound.NewDivWarSound(),
		conn:     bot.Сеть().ВебВоркер(),
		логин:    bot.Имя(),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Битва дивизий",
		СтрКонтроль_: `<span>до начала `,
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	// сам.shotTimeFull.Set(8000) // 8000 msec
	var err error
	{ // Net
		сам.net, err = div_war_net.NewDivWarNet(bot)
		if err != nil {
			return nil, fmt.Errorf("NewDivWar(): при создании DivWarNet, err=\n\t%w", err)
		}
	}
	go сам.run()
	go сам.резервТик()
	return сам, nil
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
			if сам.процВойна != nil {
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
			if сам.процВойна != nil {
				continue
			}
			сам.alarm.Уст(1)
			сам.sound.Play()
			go сам.DivWar()              // Запустить цикл непосредственного сражения
			time.Sleep(time.Second * 10) // Задержка для звука на странице
			сам.alarm.Уст(0)
		default:
			врем := сам.ВремяОстат().ПолучМилСек()
			if врем > 0 {
				continue
			}

		}
	}
}

// Ищет время до начала битвы дивизий
func (сам *DivWar) findTimeCount() {
	сам.net.Обновить()
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
	сам.ОбратВремяУст(alias.АВремя(strTime))
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
	res := сам.net.Get(linkUp)
	if res.IsErr() {
		// log._rintf("ERRO DivWar.upDivWar(): при выполнении GET-команды на подъём в атаку, err=\n\t%v\n", err)
		return
	}
	сам.СтрОбновить(res.Unwrap())
}

// Ведёт сражение
func (сам *DivWar) DivWar() {
	defer func() {
		сам.процВойна = nil
		сам.block.Unlock()
		сам.ОбратВремяУст("01")
	}()
	сам.процВойна = process_divwar.НовПроцессДивизияВойна(сам.конт) // IDivWarOn (онлайн)
	// Цикл ожидания окончания сражения
	for !сам.процВойна.ЕслиКонец().Get() {
		time.Sleep(time.Second * 1)
	}

}

// Alarm -- возвращает признак начала сражения (для браузера)
func (сам *DivWar) Alarm() ИСтатПарам {
	return сам.alarm
}
