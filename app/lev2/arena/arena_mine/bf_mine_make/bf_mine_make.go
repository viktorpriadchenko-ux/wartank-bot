// package bf_mine_make -- сделать работу в шахте
package bf_mine_make

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// ШахтаРаботать - -заставляет работать шахту
func ШахтаРаботать(конт ILocalCtx) {
	шахта := конт.Get("шахта").Val().(ИАренаШахта)
	еслиПостроено := шахта.Состояние().Получ() == cons.РежимПостроено
	еслиОжидание := шахта.Состояние().Получ() == cons.РежимОжидание
	if !(еслиПостроено || еслиОжидание) {
		return
	}
	шахтаРаботать(конт)
}
func шахтаРаботать(конт ILocalCtx) {
	сам := конт.Get("шахта").Val().(ИАренаШахта)
	сам.Сеть().Обновить()
	if err := выбратьМеталл(конт); err != nil {
		return
	}
	продукт := сам.ПродуктСейчас().Имя()
	switch продукт {
	case "руда":
		for !рудаСделать(конт) {
		}
	case "железо":
		for !железоСделать(конт) {
		}
	case "сталь":
		for !стальСделать(конт) {
		}
	case "свинец":
		for !свинецСделать(конт) {
		}
	default:
		Hassert(false, "ERRO Шахта.Сделать(): неизвестный режим производства, режим=%q\n", продукт)
	}
	сам.Состояние().Уст(cons.РежимРабота)
}

// Выбирает продукцию по возможности произвести и её количеству
func выбратьМеталл(конт ILocalCtx) error {
	сам := конт.Get("шахта").Val().(ИАренаШахта)
	var (
		диктПродукция = make(map[string]bool) // Словарь известной продукции
		lstMine       = сам.СписПолучить()
	)

	фнВыбратьПродукт := func() { // вычисляет список допустимой продукции
		диктПродукция["руда"] = true // Руда есть всегда
		диктПродукция["железо"] = false
		диктПродукция["сталь"] = false
		диктПродукция["свинец"] = false
		for _, strProd := range lstMine { // Проверить руду
			if strings.Contains(strProd, `<span class="green2">Руда</span><br/>`) {
				диктПродукция["руда"] = true
				break
			}
		}
		for _, strProd := range lstMine { // Проверить железо
			if strings.Contains(strProd, `<span class="green2">Железо</span><br/>`) {
				диктПродукция["железо"] = true
				break
			}
		}
		for _, strProd := range lstMine { // Проверить сталь
			if strings.Contains(strProd, `<span class="green2">Сталь</span><br/>`) {
				диктПродукция["сталь"] = true
				break
			}
		}
		for _, strProd := range lstMine { // Проверить свинец
			if strings.Contains(strProd, `<span class="green2">Свинец</span><br/>`) {
				диктПродукция["свинец"] = true
				break
			}
		}
		сам.ПродуктСейчас().ИмяУст("руда")
	}
	фнВыбратьПродукт()
	руда := сам.Руда().Получ()
	железо := сам.Железо().Получ()
	if диктПродукция["железо"] {
		if руда > железо*2 {
			сам.ПродуктСейчас().ИмяУст("железо")
		}
	}

	сталь := сам.Сталь().Получ()
	if диктПродукция["сталь"] {
		if железо > сталь*2 {
			сам.ПродуктСейчас().ИмяУст("сталь")
		}
	}

	свинец := сам.Свинец().Получ()
	if диктПродукция["свинец"] {
		if сталь > свинец*2 {
			сам.ПродуктСейчас().ИмяУст("свинец")
		}
		// Свинец долго делать, больше 100 не надо, а руду хоть продать можно.
		// 6 руды в час по 5 серебра = 180 серебра в сутки
		// Свинец -- 1 за 22 часа = примерно 102 серебра в сутки
		if свинец >= 100 {
			сам.ПродуктСейчас().ИмяУст("руда")
		}
	}
	return nil
}

// Создаёт руду
func рудаСделать(конт ILocalCtx) bool {
	сам := конт.Get("шахта").Val().(ИАренаШахта)
	res := сам.Сеть().Get("https://wartank.ru/production/Mine")
	if res.IsErr() {
		// log._rintf("ERRO Шахта.сделатьРуду(): при GET-команде 'начать производство руды', err=\n\t%v\n", err)
		return false
	}
	var (
		инд    int
		стрВых string
		// strTime     string
		strLink     string
		strNum      string
		еслиНайдено bool
		lstMine     = res.Unwrap()
	)

	for инд, стрВых = range lstMine {
		if strings.Contains(стрВых, `<span class="green2">Руда</span><br/>`) { // <span class="green2">Руда</span><br/>
			strNum = lstMine[инд+1]
			// strTime = lstMine[инд+3]
			strLink = lstMine[инд+10]
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return false
	}
	if !strings.Contains(strLink, `>Начать производство<`) {
		return true
	}
	// "Mine?16-1.ILinkListener-productions-0-production-startProduceLink\"><span><span>Начать производство</span></span></a>"
	// "<a class=\"simple-but border\" href=\"Mine?16-1.ILinkListener-productions-0-production-startProduceLink\"><span><span>Начать производство</span></span></a>"
	_link := strings.TrimPrefix(strLink, `<a class="simple-but border" href="`)
	_link = strings.TrimSuffix(_link, "\"><span><span>Начать производство</span></span></a>")
	strLink = "https://wartank.ru/production/" + _link
	// https://wartank.ru/production/Mine?19-1.ILinkListener-productions-0-production-startProduceLink
	res = сам.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Шахта.сделатьРуду(): при GET-команде 'начать производство руды', err=\n\t%v\n", err)
		return false
	}
	lstMine = res.Unwrap()
	for _, стрВых = range lstMine {
		if strings.Contains(стрВых, `><span><span>Начать производство</span></span></a>`) {
			return false
		}
	}
	// сам.СтрОбновить(lstMine)
	//сам.ОбратВремяУст(АВремя(strTime))
	lstNum := strings.Split(strNum, `Кол-во: <span class="green2">`)
	strNum = lstNum[1]
	lstNum = strings.Split(strNum, `</span><br/>`)
	strNum = lstNum[0]
	iNum, err := strconv.Atoi(strNum)
	if err != nil {
		// log._rintf("ERRO Шахта.сделатьРуду(): кол-во(%v) не число, err=\n\t%v\n", strNum, err)
		return false
	}
	сам.ПродуктСейчас().Уст(iNum)
	сам.ПродуктСейчас().ИмяУст("руда")
	return true
}

// Создаёт железо
func железоСделать(конт ILocalCtx) bool {
	сам := конт.Get("шахта").Val().(ИАренаШахта)
	var (
		lstMine = сам.СписПолучить()
		ind     int
		strOut  string
		// strTime     string
		strLink     string
		strNum      string
		еслиНайдено bool
	)
	for ind, strOut = range lstMine {
		if strings.Contains(strOut, `<span class="green2">Железо</span><br/>`) {
			// <span class="green2">Железо</span><br/>
			strNum = lstMine[ind+1]
			// Кол-во: <span class="green2">1</span><br/>
			// strTime = lstMine[ind+3]
			// <a class="simple-but border" href="Mine?4-1.ILinkListener-productions-1-production-startProduceLink"><span><span>Начать производство</span></span></a>
			strLink = lstMine[ind+10]
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return true
	}
	lstLink := strings.Split(strLink, `<a class="simple-but border" href="`)
	strLink = lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
	strLink = "https://wartank.ru/production/" + lstLink[0]
	// https://wartank.ru/production/Mine?4-1.ILinkListener-productions-1-production-startProduceLink
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO MineNet.makeFerrum(): при GET-команде 'начать производство железа', err=\n\t%v\n", err)
		return false
	}
	lstMine = res.Unwrap()
	for _, strOut := range lstMine { // Проверка на базу
		if strings.Contains(strOut, `<title>База</title>`) {
			// log._rintf("ERRO MineNet.makeFerrum(): при обновлении lstMine найден lstBase")
			return false
		}
	}
	сам.СтрОбновить(lstMine)
	// сам.ОбратВремяУст(АВремя(strTime))
	lstNum := strings.Split(strNum, `Кол-во: <span class="green2">`)
	strNum = lstNum[1]
	lstNum = strings.Split(strNum, `</span><br/>`)
	strNum = lstNum[0]
	iNum, err := strconv.Atoi(strNum)
	if err != nil {
		// log._rintf("ERRO MineNet.makeFerrum(): кол-во(%v) не число, err=\n\t%v\n", strNum, err)
		return false
	}
	сам.ПродуктСейчас().Уст(iNum)
	сам.ПродуктСейчас().ИмяУст("железо")
	return true
}

// Создаёт сталь
func стальСделать(конт ILocalCtx) bool {
	сам := конт.Get("шахта").Val().(ИАренаШахта)
	var (
		lstMine = сам.СписПолучить()
		ind     int
		strOut  string
		// strTime     string
		strLink     string
		strNum      string
		еслиНайдено bool
	)
	for ind, strOut = range lstMine {
		if strings.Contains(strOut, `<span class="green2">Сталь</span><br/>`) {
			strNum = lstMine[ind+1]
			// strTime = lstMine[ind+3]
			strLink = lstMine[ind+10]
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return true
	}
	lstLink := strings.Split(strLink, `<a class="simple-but border" href="`)
	strLink = lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
	strLink = "https://wartank.ru/production/" + lstLink[0]
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO MineNet.makeSteel(): при GET-команде 'начать производство стали', err=\n\t%v\n", err)
		return false
	}
	lstMine = res.Unwrap()
	for _, strOut := range lstMine { // Проверка на базу
		if strings.Contains(strOut, `<title>База</title>`) {
			// log._rintf("ERRO MineNet.makeSteel(): при обновлении lstMine найден lstBase")
			return false
		}
	}
	сам.СтрОбновить(lstMine)
	// сам.ОбратВремяУст(АВремя(strTime))
	lstNum := strings.Split(strNum, `Кол-во: <span class="green2">`)
	strNum = lstNum[1]
	lstNum = strings.Split(strNum, `</span><br/>`)
	strNum = lstNum[0]
	iNum, err := strconv.Atoi(strNum)
	if err != nil {
		// log._rintf("ERRO MineNet.makeSteel(): кол-во(%v) не число, err=\n\t%v\n", strNum, err)
		return false
	}
	сам.ПродуктСейчас().Уст(iNum)
	сам.ПродуктСейчас().ИмяУст("сталь")
	return true
}

// Создаёт свинец
func свинецСделать(конт ILocalCtx) bool {
	сам := конт.Get("шахта").Val().(ИАренаШахта)
	var (
		lstMine = сам.СписПолучить()
		ind     int
		strOut  string
		// strTime     string
		strLink     string
		strNum      string
		еслиНайдено bool
	)
	for ind, strOut = range lstMine {
		if strings.Contains(strOut, `<span class="green2">Свинец</span><br/>`) {
			strNum = lstMine[ind+1]
			// strTime = lstMine[ind+3]
			strLink = lstMine[ind+10]
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return true
	}
	lstLink := strings.Split(strLink, `<a class="simple-but border" href="`)
	strLink = lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
	strLink = "https://wartank.ru/production/" + lstLink[0]
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Шахта.сделатьСвинец(): при GET-команде 'начать производство стали', err=\n\t%v\n", err)
		return false
	}
	lstMine = res.Unwrap()
	for _, strOut := range lstMine { // Проверка на базу
		if strings.Contains(strOut, `<title>База</title>`) {
			// log._rintf("ERRO Шахта.сделатьСвинец(): при обновлении lstMine найден lstBase")
			return false
		}
	}
	сам.СтрОбновить(lstMine)
	// сам.ОбратВремяУст(АВремя(strTime))
	lstNum := strings.Split(strNum, `Кол-во: <span class="green2">`)
	strNum = lstNum[1]
	lstNum = strings.Split(strNum, `</span><br/>`)
	strNum = lstNum[0]
	iNum, err := strconv.Atoi(strNum)
	if err != nil {
		// log._rintf("ERRO Шахта.сделатьСвинец(): кол-во(%v) не число, err=\n\t%v\n", strNum, err)
		return false
	}
	сам.ПродуктСейчас().Уст(iNum)
	сам.ПродуктСейчас().ИмяУст("свинец")
	return true
}
