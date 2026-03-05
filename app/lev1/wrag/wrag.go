package wrag

import (
	"strconv"
	"strings"
	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Танк враг на битве, сражении, дуэли, войне
*/

// Враг -- объект врага
type Враг struct {
	сервер   ИПриложение
	лог      ILogBuf
	здоровье int //
}

// НовВраг -- возвращает новый объект врага
func НовВраг(конт IKernelCtx, app ИПриложение, lstBattle []string) *Враг {
	лог := NewLogBuf()
	лог.Debug("НовВраг()\n")
	сам := &Враг{
		сервер: app,
		лог:    лог,
	}
	сам.update(lstBattle)
	return сам
}

// Обновляет вражину
func (сам *Враг) update(lstBattleOn []string) {
	// <img class="tank-img" src="/tankimg?c=2&amp;k=1&amp;m=0-2,1-2,2-0,3-2,5-2,6-0&amp;t=png" alt="Тень Брата">
	var (
		ind         int
		strOut      string
		еслиНайдено bool
	)
	for ind, strOut = range lstBattleOn {
		if strings.Contains(strOut, `<img class="tank-img" src="/`) {
			// Убедиться, что это не свой танк
			if strings.Contains(strOut, "prospero tank") {
				continue
			}
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Не нашёл метку врага
		сам.лог.Err("update(): не нашёл метку врага\n")
		сам.здоровье = 800
		return
	}
	// Вражина найдена, ищем настоящее здоровье
	ind += 13
	strOut = lstBattleOn[ind]
	lstHealth := strings.Split(strOut, `<span>`)
	strHealth := lstHealth[1]
	lstHealth = strings.Split(strHealth, `</span>`)
	strHealth = lstHealth[0]
	iHealth, err := strconv.Atoi(strHealth)
	if err != nil {
		сам.лог.Err("update(): здоровье(%v) не число, err=\n\t%v\n", strHealth, err)
		сам.здоровье = 800
		return
	}
	сам.здоровье = iHealth
}
