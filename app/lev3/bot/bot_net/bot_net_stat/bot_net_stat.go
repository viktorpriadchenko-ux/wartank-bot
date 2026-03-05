// package bot_net_stat -- сетевая статистика бота
package bot_net_stat

import (
	"strconv"
	"strings"

	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Исходник предоставляет тип для парсинга статистики уровня бота.
*/

// БотНетСтат -- статистика уровня бота
type БотСетьСтат struct {
	прилож ИПриложение
	бот    ИБот
	лог    ILogBuf
	конт   ILocalCtx
	ангар  ИАренаАнгар
}

// НовБотСетьСтат -- возвращает новый *NetStat
func НовБотСетьСтат(конт ILocalCtx) *БотСетьСтат {
	лог := NewLogBuf()
	лог.Info("НовБотСетьСтат()\n")
	bot := конт.Get("бот").(ИБот)
	сам := &БотСетьСтат{
		прилож: конт.Get("сервер").(ИПриложение),
		бот:    bot,
		лог:    лог,
		конт:   конт,
		ангар:  конт.Get("ангар").Val().(ИАренаАнгар),
	}
	return сам
}

func (сам *БотСетьСтат) Update() {
	// _mt.Println("\n===NetStat.Update()===")
	сам.findLevelTank()
	сам.findLevelProgress()
	сам.атакаНайти()
	сам.броняНайти()
	сам.точностьНайти()
	сам.прочностьНайти()
	сам.мощностьНайти()
	сам.игроковОнлайнНайти()
}

// Ищет в теле текста ангара мощность танка
func (сам *БотСетьСтат) мощностьНайти() {
	lstAngar := сам.ангар.СписПолучить()
	if len(lstAngar) == 0 {
		// log._rintf("WARN NetStat.findPower(): lstAngar is empty\n")
		return
	}
	var strOut string
	for _, strPower := range lstAngar {
		if strings.Contains(strPower, `/images/icons/power.png?2`) {
			strOut = strPower
			break
		}
	}
	// Выделить мощность
	lstPower := strings.Split(strOut, `<img src="/images/icons/power.png?2" height="14" width="14"> <span class="green2">Танковая мощь: `)
	strPower := lstPower[1]
	lstPower = strings.Split(strPower, `</span>`)
	strPower = lstPower[0]
	iPower, ош := strconv.Atoi(strPower)
	Hassert(ош == nil, "NetStat.мощностьНайти(): мощность(%v) не число, ош=\n\t%v\n", strPower, ош)
	сам.бот.Стата().Мощь().Уст(iPower)
}

// Ищет в теле текста ангара прочность танка
func (сам *БотСетьСтат) прочностьНайти() {
	var (
		strOut      string
		lstAngar    = сам.ангар.СписПолучить()
		еслиНайдено bool
	)
	if len(lstAngar) == 0 {
		// log._rintf("WARN NetStat.findHard(): lstAngar пустой\n")
		return
	}
	for _, strOut = range lstAngar {
		if strings.Contains(strOut, `/images/icons/durability.png?1`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// Выделить прочность
	lstHard := strings.Split(strOut, `<img width="14" height="14" src="/images/icons/durability.png?1" alt="Прочность" title="Прочность"/> Прочность <span class="green2">`)
	strHard := lstHard[1]
	lstHard = strings.Split(strHard, `</span><br/>`)
	strHard = lstHard[0]
	iHard, ош := strconv.Atoi(strHard)
	Hassert(ош == nil, "NetStat.прочностьНайти(): прочность(%v) не число, ош=\n\t%v\n", strHard, ош)
	сам.бот.Стата().Прочность().Уст(iHard)
}

// Ищет в теле текста ангара точность танка
func (сам *БотСетьСтат) точностьНайти() {
	var (
		strOut      string
		lstAngar    = сам.ангар.СписПолучить()
		еслиНайдено bool
	)
	if len(lstAngar) == 0 {
		// log._rintf("WARN NetStat.findFyne(): lstAngar пустой\n")
		return
	}
	for _, strOut = range lstAngar {
		if strings.Contains(strOut, `/images/icons/accuracy.png?1`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// Выделить броню
	lstFyne := strings.Split(strOut, `<img width="14" height="14" src="/images/icons/accuracy.png?1" alt="Точность" title="Точность"/> Точность <span class="green2">`)
	strFyne := lstFyne[1]
	lstFyne = strings.Split(strFyne, `</span><br/>`)
	strFyne = lstFyne[0]
	iFyne, ош := strconv.Atoi(strFyne)
	Hassert(ош == nil, "NetStat.точностьНайти(): точность(%v) не число, ош=\n\t%v\n", strFyne, ош)
	сам.бот.Стата().Точность().Уст(iFyne)
}

// Ищет в теле текста ангара броню танка
func (сам *БотСетьСтат) броняНайти() {
	var (
		strOut      string
		lstAngar    = сам.ангар.СписПолучить()
		еслиНайдено bool
	)
	if len(lstAngar) == 0 {
		return
	}
	for _, strOut = range lstAngar {
		if strings.Contains(strOut, `/images/icons/armor.png?1`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// Выделить броню
	lstArmor := strings.Split(strOut, `<img width="14" height="14" src="/images/icons/armor.png?1" alt="Броня" title="Броня"/> Броня <span class="green2">`)
	strArmor := lstArmor[1]
	lstArmor = strings.Split(strArmor, `</span><br/>`)
	strArmor = lstArmor[0]
	iArmor, ош := strconv.Atoi(strArmor)
	Hassert(ош == nil, "NetStat.броняНайти(): броня(%v) не число, ош=\n\t%v\n", strArmor, ош)
	сам.бот.Стата().Броня().Уст(iArmor)
}

// Ищет в теле текста ангара уровень танка
func (сам *БотСетьСтат) findLevelTank() {
	lstAngar := сам.ангар.СписПолучить()
	if len(lstAngar) == 0 {
		// log._rintf("ERRO NetStat.findLevelTank(): пустой lstAngar")
		return
	}
	var strOut string
	for _, strLevel := range lstAngar {
		if strings.Contains(strLevel, "/images/icons/level.png?2") {
			strOut = strLevel
			break
		}
	}
	// Выделить уровень
	lstLevel := strings.Split(strOut, `<td><div class="value-block lh1"><span><span><img height="14" width="14" src="/images/icons/level.png?2"/> `)
	strLevel := lstLevel[1]
	lstLevel = strings.Split(strLevel, `</span></span></div></td>`)
	strLevel = lstLevel[0]
	iLevel, ош := strconv.Atoi(strLevel)
	Hassert(ош == nil, "NetStat.прочностьНайти(): уровень(%v) не число, ош=\n\t%v\n", strLevel, ош)
	сам.ангар.Уровень().Уст(iLevel)
}

// Ищет в теле текста ангара прогресс уровня танка танка
func (сам *БотСетьСтат) findLevelProgress() {
	lstAngar := сам.ангар.СписПолучить()
	var strOut string
	for _, strProg := range lstAngar {
		if strings.Contains(strProg, `class="progr"`) {
			strOut = strProg
			break
		}
	}
	// Выделить прогресс
	lstProg := strings.Split(strOut, `<td class="progr"><div class="scale-block"><div class="scale" style="width:`)
	strProg := lstProg[1]
	lstProg = strings.Split(strProg, `%;">&nbsp;</div></div></td>`)
	strProg = lstProg[0]
	iProg, ош := strconv.Atoi(strProg)
	Hassert(ош == nil, "NetStat.прогрессНайти(): прогресс(%v) не число, ош=\n\t%v\n", strProg, ош)
	сам.ангар.Прогресс().Уст(iProg)
}

// Ищет в теле текста ангара силу атаки танка
func (сам *БотСетьСтат) атакаНайти() {
	var (
		strOut      string
		lstAngar    = сам.ангар.СписПолучить()
		еслиНайдено bool
	)
	for _, strOut = range lstAngar {
		if strings.Contains(strOut, `/images/icons/attack.png?1`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// Выделить атаку
	списАтака := strings.Split(strOut, `<img width="14" height="14" src="/images/icons/attack.png?1" alt="Атака" title="Атака"/> Атака <span class="green2">`)
	стрАтака := списАтака[1]
	списАтака = strings.Split(стрАтака, `</span><br/>`)
	стрАтака = списАтака[0]
	целАтака, ош := strconv.Atoi(стрАтака)
	Hassert(ош != nil, "NetStat.атакаНайти(): атака(%v) не число, ош=\n\t%v\n", стрАтака, ош)
	сам.бот.Стата().Атака().Уст(целАтака)
}

// Ищет в теле текста ангара силу атаки танка
func (сам *БотСетьСтат) игроковОнлайнНайти() {
	lstAngar := сам.ангар.СписПолучить()
	var strOut string
	for _, стрАтака := range lstAngar {
		if strings.Contains(стрАтака, `>Онлайн</a>: `) {
			strOut = стрАтака
			break
		}
	}
	// Выделить число игроков онлайн
	lstAngar = strings.Split(strOut, `<span class="yellow1">`)
	if len(lstAngar) <= 1 {
		сам.конт.Cancel()
		return
	}
	strOnline := lstAngar[1]
	lstAngar = strings.Split(strOnline, `</span>`)
	strOnline = lstAngar[0]
	iOnline, ош := strconv.Atoi(strOnline)

	Hassert(ош == nil, "NetStat.findOnline(): игроков онлайн(%v) не число, ош=\n\t%v\n", strOnline, ош)
	сам.ангар.ИгрокиОнлайн().Уст(iOnline)
}
