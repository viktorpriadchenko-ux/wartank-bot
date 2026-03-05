package arena_polygon

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_polygon/bf_polygon_build"
	"wartank/app/lev2/arena/arena_polygon/bf_polygon_level"
	"wartank/app/lev2/arena/arena_polygon/bf_polygon_make"
	"wartank/app/lev2/arena/arena_polygon/bf_polygon_upgrade"
	"wartank/app/lev2/arena/arena_polygon/bf_polygon_upgrade_fast"
)

/*
	Объект полигона на базе
*/

const (
	времОжидПлат    = "05:00" // Время ожидания платного ускорения
	времОжидБесплат = "30:00" // Время ожидания бесплатного ускорения
	стрПрочность    = "прочность"
	стрТочность     = "точность"
	стрБроня        = "броня"
	стрАтака        = "атака"
)

// АренаПолигон -- объект полигона на базе
type АренаПолигон struct {
	ИАренаСтроение
	конт     ILocalCtx
	танкСтат ИТанкСтат
	лог      ILogBuf
}

// НовПолигон -- возвращает новый *Polygon
func НовПолигон(конт ILocalCtx) *АренаПолигон {
	лог := NewLogBuf()
	бот := конт.Get("бот").Val().(ИБот)
	лог.Info("НовПолигон(): бот=%s\n", бот.Имя())
	сам := &АренаПолигон{
		танкСтат: бот.Стата(),
		лог:      лог,
		конт:     конт,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Полигон",
		СтрКонтроль_: `<title>Полигон</title>`,
		СтрУрл_:      "https://wartank.ru/polygon",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	конт.Set("полигон", сам, "Полигон бота")
	_ = ИАренаПолигон(сам)
	return сам
}

func (сам *АренаПолигон) Пуск() {
	сам.ИАренаСтроение.Пуск()
	bf_polygon_build.ПолигонПостроить(сам.конт)
	bf_polygon_upgrade.ПолигонАпгрейд(сам.конт)
	bf_polygon_upgrade_fast.ПолигонАпгрейдБесплатно(сам.конт)
	bf_polygon_make.ПолигонВключить(сам.конт)
	bf_polygon_level.ПолигонУровень(сам.конт)
	bf_polygon_upgrade.ПолигонАпгрейд(сам.конт)
	сам.времяОбнов()
	сам.усилениеПровер()
}

// Обновляет оставшееся время полигона
//
//	Этот объект сам описывает своё время
func (сам *АренаПолигон) времяОбнов() {
	var (
		strLastTime string
		еслиНайдено bool
		isSet       bool
		lstPolygon  = сам.СписПолучить()
	)
	defer func() {
		if !isSet {
			сам.ОбратВремяУст("05")
		}
	}()
	for _, lastTime := range lstPolygon {
		if strings.Contains(lastTime, `>Осталось: `) {
			strLastTime = lastTime
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Время полигона вышло
		return
	}
	lstTime := strings.Split(strLastTime, `>Осталось: `)
	strLastTime = lstTime[1]
	lstTime = strings.Split(strLastTime, `</span>`)
	strLastTime = lstTime[0]
	сам.ОбратВремяУст(АВремя(strLastTime))
	isSet = true
}

// Проверяет что именно активировано
func (сам *АренаПолигон) усилениеПровер() {
	var (
		еслиНайдено bool
		lstPolygon  = сам.СписПолучить()
		ind         = 0
		strOut      string
	)

	for ind, strOut = range lstPolygon {
		if strings.Contains(strOut, `<span>Активно</span>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	if ind < 9 {
		return // Недостаточно строк перед маркером
	}
	ind -= 9
	strOut = lstPolygon[ind]
	форсажИмя := ""
	switch { // Вычисляем контрольную строку
	case strings.Contains(strOut, `>улучшение точности<`):
		форсажИмя = стрТочность
	case strings.Contains(strOut, `>увеличение прочности<`):
		форсажИмя = стрПрочность
	case strings.Contains(strOut, `>усиление брони<`):
		форсажИмя = стрБроня
	case strings.Contains(strOut, `>усиление атаки<`):
		форсажИмя = стрАтака
	}
	// Вычислим на сколько
	if ind+1 >= len(lstPolygon) {
		return
	}
	strOut = lstPolygon[ind+1]
	lstOut := strings.Split(strOut, `<span class="green2">+`)
	if len(lstOut) < 2 {
		return
	}
	strOut = lstOut[1]
	lstOut = strings.Split(strOut, ` на `)
	strOut = lstOut[0]
	iForce, err := strconv.Atoi(strOut)
	if err != nil {
		// log._rintf("NetPolygon.checkTime(): force(%v) not number, err=\n\t%v\n", strOut, err)
		return
	}
	сам.танкСтат.ФорсажОбнов(форсажИмя, iForce)
	сам.ПродуктСейчас().ИмяУст("усиление-" + форсажИмя)
	сам.ПродуктСейчас().Уст(iForce)
}
