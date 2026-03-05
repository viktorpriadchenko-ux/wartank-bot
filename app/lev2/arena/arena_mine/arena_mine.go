// package arena_mine -- объект шахты на базе
package arena_mine

import (
	"fmt"
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	// . "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev1"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_mine/bf_mine_accelerate"
	"wartank/app/lev2/arena/arena_mine/bf_mine_build"
	"wartank/app/lev2/arena/arena_mine/bf_mine_make"
	"wartank/app/lev2/arena/arena_mine/bf_mine_take"
	"wartank/app/lev2/arena/arena_mine/bf_mine_time_work"
)

// АренаШахта -- объект шахты на базе
type АренаШахта struct {
	ИАренаСтроение
	лог    ИВебЛог
	база   ИАренаБаза
	руда   ИСтатПарам
	железо ИСтатПарам
	сталь  ИСтатПарам
	свинец ИСтатПарам
	конт   ILocalCtx
}

// НовШахта -- возвращает новый *Mine
func НовШахта(конт ILocalCtx) *АренаШахта {
	сам := &АренаШахта{
		конт:   конт,
		база:   конт.Get("база").Val().(ИАренаБаза),
		руда:   lev1.НовСтатПарам("руда"),
		железо: lev1.НовСтатПарам("железо"),
		сталь:  lev1.НовСтатПарам("сталь"),
		свинец: lev1.НовСтатПарам("свинец"),
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        сам.конт,
		АренаИмя_:    "Шахта",
		СтрКонтроль_: `<span class="green2">Руда</span><br/>`,
		СтрУрл_:      "https://wartank.ru/production/Mine",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	сам.лог = сам.ВебЛог()
	go сам.пуск()
	сам.лог.ОтклВывод()
	сам.лог.Добавить("Шахта.НовШахта(): бот=%q\n", конт.Get("бот").Val().(ИБот).Имя())
	конт.Set("шахта", сам, "Шахта бота")
	_ = ИАренаШахта(сам)
	return сам
}

func (сам *АренаШахта) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_mine_build.ШахтаПостроить(сам.конт)
	bf_mine_accelerate.ШахтаУскорить(сам.конт)
	bf_mine_take.ШахтаЗабрать(сам.конт)
	bf_mine_make.ШахтаРаботать(сам.конт)
	bf_mine_time_work.ШахтаРаботаВремя(сам.конт)
}

// пуск -- запускает обработку шахты
func (сам *АренаШахта) пуск() {}

// Проверяет количество продукта в шахте
func (сам *АренаШахта) количествоПолучить() (bool, error) {
	сам.лог.Добавить("Шахта.количествоПолучить()\n")
	var (
		strOut      string
		еслиНайдено bool
		режим       string
	)

	lstMine := сам.Сеть().ВебВоркер().Получ("https://wartank.ru/buildings")

	/*
		Режим (руда-1):
		<td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/ore.png?2" alt="ore"/>&nbsp;1</div></td>

		Время (+8 строк):
		<td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/ore.png?2" alt="ore"/>&nbsp;1</div></td>
	*/
	for _, strOut = range lstMine {
		// Руда текущее
		if strings.Contains(strOut, `src="/images/icons/ore.png?2" alt="ore"`) {
			// <td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/ore.png?2" alt="ore"/>&nbsp;1</div></td>
			еслиНайдено = true
			режим = "руда"
			break
		}
		// Железо текущее
		if strings.Contains(strOut, `src="/images/icons/iron.png?2" alt="iron"`) {
			// <td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/iron.png?2" alt="iron"/>&nbsp;2</div></td>
			еслиНайдено = true
			режим = "железо"
			break
		}
		// Сталь текущее
		if strings.Contains(strOut, `src="/images/icons/steel.png?2" alt="steel"`) {
			// <td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/steel.png?2" alt="iron"/>&nbsp;2</div></td>
			еслиНайдено = true
			режим = "сталь"
			break
		}
		// Свинец текущее
		if strings.Contains(strOut, `src="/images/icons/plumbum.png?2" alt="plumbum"`) {
			// <td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/plumbum.png?2" alt="iron"/>&nbsp;2</div></td>
			еслиНайдено = true
			режим = "свинец"
			break
		}
	}
	if !еслиНайдено {
		сам.лог.Добавить("Шахта.количествоПолучить(): не надо\n")
		return true, nil
	}
	switch режим {
	case "руда":
		_число := strings.TrimPrefix(strOut, `<td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/ore.png?2" alt="ore"/>&nbsp;`)
		_число = strings.TrimSuffix(_число, `</div></td>`)
		iNum, err := strconv.Atoi(_число)
		if err != nil {
			сам.лог.Добавить("ОШИБКА Шахта.количествоПолучить(): кол-во руды (%v) не число, err=\n\t%v\n", _число, err)
			return false, fmt.Errorf("")
		}
		сам.ПродуктСейчас().Уст(iNum)
		сам.ПродуктСейчас().ИмяУст("руда")
		сам.лог.Добавить("Шахта.количествоПолучить(): кол-во руды = %v\n", iNum)
	case "железо":
		_число := strings.TrimPrefix(strOut, `<td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/iron.png?2" alt="iron"/>&nbsp;`)
		_число = strings.TrimSuffix(_число, `</div></td>`)
		iNum, err := strconv.Atoi(_число)
		if err != nil {
			сам.лог.Добавить("ОШИБКА Шахта.количествоПолучить(): кол-во железа (%v) не число, err=\n\t%v\n", _число, err)
			return false, fmt.Errorf("")
		}
		сам.ПродуктСейчас().Уст(iNum)
		сам.ПродуктСейчас().ИмяУст("железо")
		сам.лог.Добавить("Шахта.количествоПолучить(): кол-во железа = %v\n", iNum)
	case "сталь":
		_число := strings.TrimPrefix(strOut, `<td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/steel.png?2" alt="steel"/>&nbsp;`)
		_число = strings.TrimSuffix(_число, `</div></td>`)
		iNum, err := strconv.Atoi(_число)
		if err != nil {
			сам.лог.Добавить("ОШИБКА Шахта.количествоПолучить(): кол-во стали (%v) не число, err=\n\t%v\n", _число, err)
			return false, fmt.Errorf("")
		}
		сам.ПродуктСейчас().Уст(iNum)
		сам.ПродуктСейчас().ИмяУст("сталь")
		сам.лог.Добавить("Шахта.количествоПолучить(): кол-во стали = %v\n", iNum)
	case "свинец":
		_число := strings.TrimPrefix(strOut, `<td class="vam"><div class="nwr pr5 gray1"><img class="rico vm" src="/images/icons/plumbum.png?2" alt="plumbum"/>&nbsp;`)
		_число = strings.TrimSuffix(_число, `</div></td>`)
		iNum, err := strconv.Atoi(_число)
		if err != nil {
			сам.лог.Добавить("ОШИБКА Шахта.количествоПолучить(): кол-во свинца (%v) не число, err=\n\t%v\n", _число, err)
			return false, fmt.Errorf("")
		}
		сам.ПродуктСейчас().Уст(iNum)
		сам.ПродуктСейчас().ИмяУст("свинец")
		сам.лог.Добавить("Шахта.количествоПолучить(): кол-во свинца = %v\n", iNum)
	default:
		сам.лог.Добавить("Шахта.количествоПолучить(): неизвестный режим (%v)\n", режим)
		return false, fmt.Errorf("")
	}
	return true, nil
}

// // Проверяет ускорение строительства FIXME: не работает
// func (сам *АренаШахта) ускорениеПровер() {
// 	сам.лог.Добавить("")
// 	списСтр := сам.Сеть().ВебВоркер().Получ("http://wartank.ru/buildings")
// 	// <span class="green2">Шахта - 0</span><br/>
// 	var (
// 		еслиНайти bool
// 		стр       string
// 	)
// 	for _, стр = range списСтр {
// 		if strings.Contains(стр, `<span class="green2">Шахта - `) {
// 			еслиНайти = true
// 			break
// 		}
// 	}
// 	if !еслиНайти {
// 		сам.лог.Добавить("Шахта.ускорениеПровер(): не надо\n")
// 		return
// 	}
// 	сам.лог.Добавить("Шахта.ускорениеПровер(): надо\n")
// }

// Обновляет текущий уровень шахты (может быть не построена)
func (сам *АренаШахта) уровеньОбновить() bool {
	сам.лог.Добавить("Шахта.уровеньОбновить()\n")
	списСтр := сам.Сеть().ВебВоркер().Получ("http://wartank.ru/buildings")

	// <span class="green2">Шахта - 0</span><br/>
	var (
		еслиНайти = false
		стр       = ""
	)
	for _, стр = range списСтр {
		if strings.Contains(стр, `<span class="green2">Шахта - `) {
			еслиНайти = true
			break
		}
	}
	if !еслиНайти {
		сам.лог.Добавить("Шахта.уровеньОбновить(): нет уровня\n")
		return false
	}
	_стр := strings.TrimPrefix(стр, `<span class="green2">Шахта - `)
	_стр = strings.TrimSuffix(_стр, `</span><br/>`)
	иУровень, ош := strconv.Atoi(_стр)
	if ош != nil {
		сам.лог.Добавить("ОШИБКА Шахта.уровеньОбновить(): строка уровня сбойная, стр=%q, ош=\n\t%v\n", стр, ош)
		return false
	}
	сам.Уровень().Уст(иУровень)
	сам.лог.Добавить("Шахта.уровеньОбновить(): уровень=%v\n", иУровень)
	return true
}

// Сделать -- вызывается с базы, если она обнаружила, что пора сделать продукцию
func (сам *АренаШахта) Сделать() {

}

// Свинец -- возвращает объект свинца
func (сам *АренаШахта) Свинец() ИСтатПарам {
	return сам.свинец
}

// Сталь -- возвращает объект стали
func (сам *АренаШахта) Сталь() ИСтатПарам {
	return сам.сталь
}

// Железо -- возвращает объект железа
func (сам *АренаШахта) Железо() ИСтатПарам {
	return сам.железо
}

// Руда -- возвращает объект руды
func (сам *АренаШахта) Руда() ИСтатПарам {
	return сам.руда
}
