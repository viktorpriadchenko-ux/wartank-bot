// package death_wait -- заставляет ожидать начало битвы
package death_wait

import (
	"strings"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// СражениеОжидание -- ожидатель начала битвы
type СражениеОжидание struct {
	ИАренаСтроение
	конт ILocalCtx
}

// НовСражениеОжидание -- возвращает новый ожидатель битвы
func НовСражениеОжидание(конт IKernelCtx) *СражениеОжидание {
	сам := &СражениеОжидание{
		конт: конт,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Ожидание сражения",
		СтрКонтроль_: `<title>Сражения</title>`,
		СтрУрл_:      "https://wartank.ru/pve",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	return сам
}

// Ожидать -- ожидает начало сражения
func (сам *СражениеОжидание) Ожидать() {

	// Зайти в цикл ожидания сражения
	for {
		// countTime := сам.ВремяОпрос().Получ()
		// if countTime > 0 {
		// 	time.Sleep(time.Millisecond * 500)
		// 	// log.Printf("BattleWait.Wait(): countTime=%v\n", сам.CountDown().String())
		// 	continue
		// }
		стрВрем := сам.ждать()
		if стрВрем == "" {
			return
		}
		лстВрем := strings.Split(стрВрем, ":")
		стрЧас := лстВрем[0]
		if стрЧас > "00" {
			time.Sleep(time.Hour * 1)
			continue
		}
		стрМин := лстВрем[1]
		if стрМин > "10" {
			time.Sleep(time.Minute * 10)
			continue
		}
		if стрМин > "01" {
			time.Sleep(time.Minute * 1)
			continue
		}
		if "00:00:05" < стрВрем && стрВрем < "00:00:59" {
			time.Sleep(time.Second * 5)
			continue
		}
		time.Sleep(time.Second * 1)
	}
}

// Ждёт пока время не обнулится
func (сам *СражениеОжидание) ждать() string {
	сам.Обновить()
	var (
		strOut      string
		lstBattle   = сам.СписПолучить()
		еслиНайдено bool
	)
	for _, strOut = range lstBattle {
		if strings.Contains(strOut, `<span>до начала `) {
			еслиНайдено = true
			break
		}
		// if strings.Contains(strOut, `>ОБЫЧНЫЕ<`) { // Это уже битва
		// 	if len(сам.chBattle) == 0 {
		// 		сам.chBattle <- 1
		// 	}
		// 	return
		// }
	}
	if !еслиНайдено { // Сражение уже идёт
		return ""
	}
	// Найдена строка ожидания начала сражения
	lstTime := strings.Split(strOut, `<span>до начала `)
	strTime := lstTime[1]
	lstTime = strings.Split(strTime, ` (`)
	strTime = lstTime[0]
	сам.ПродуктВремяСейчас().Set(strTime)
	return strTime
}
