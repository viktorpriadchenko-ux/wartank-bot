// package bf_polygon_make -- бизнес-функция активации усиления полигона
package bf_polygon_make

import (
	"log"
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

const (
	времОжидПлат    = "05:00" // Время ожидания платного ускорения
	времОжидБесплат = "30:00" // Время ожидания бесплатного ускорения
	стрПрочность    = "прочность"
	стрТочность     = "точность"
	стрБроня        = "броня"
	стрАтака        = "атака"
)

// ПолигонВключить -- активация усиления полигона
func ПолигонВключить(конт ILocalCtx) {
	полигон_ := конт.Get("полигон")
	if полигон_ == nil { // Может быть ещё не построен
		return
	}
	полигон := полигон_.Val().(ИАренаПолигон)
	еслиОжидание := полигон.Состояние().Получ() == cons.РежимОжидание
	еслиПостроено := полигон.Состояние().Получ() == cons.РежимПостроено
	if !(еслиОжидание || еслиПостроено) {
		return
	}
	усилениеДобавить(конт)
}

// Выбирает самый слабый параметр и усиливает его
func усилениеДобавить(конт ILocalCtx) {
	полигон := конт.Get("полигон").Val().(ИАренаПолигон)

	lstPolygon := полигон.СписПолучить()
	if len(lstPolygon) == 0 {
		полигон.Обновить()
		lstPolygon = полигон.СписПолучить()
	}
	танкСтат := конт.Get("танкСтат").Val().(ИТанкСтат)
	цАтака := танкСтат.Атака().Получ()
	цБроня := танкСтат.Броня().Получ()
	цТочность := танкСтат.Точность().Получ()
	цПрочность := танкСтат.Прочность().Получ()

	strParam := стрАтака
	iParam := цПрочность
	{
		/*
			Вычислить самый слабый параметр.

			Политика вычислений:
				1) hard -- прочность, самый низкоприоритетный параметр
				2) armor -- броня, чуть лучше power
				3) fyne -- точность, чуть лучше armor
				4) attack -- атака, самый важный
		*/

		if цБроня <= iParam {
			iParam = цБроня
			strParam = стрБроня
		}
		if цТочность <= iParam {
			iParam = цТочность
			strParam = стрТочность
		}
		if цАтака < iParam {
			strParam = стрАтака
		}
	}

	// Найти нужную строку активации
	var (
		ind    int
		стрРез string
	)
	switch strParam {
	case стрАтака: // Усиливаем атаку
		фнУсиление := func() bool {
			// <a class="simple-but border" href="polygon?572-1.ILinkListener-buffs-0-buff-getFreeLink"><span><span>Получить бесплатно</span></span></a>
			for ind, стрРез = range lstPolygon {
				if strings.Contains(стрРез, `>усиление атаки<`) {
					ind += 8
					if ind >= len(lstPolygon) {
						return false
					}
					стрРез = lstPolygon[ind]
					if стрРез == "" {
						return false
					}
					if strings.Contains(стрРез, "Активировать за") {
						return false
					}
					return true
				}
			}
			return false
		}
		полигон.Обновить()
		lstPolygon = полигон.СписПолучить()
		if фнУсиление() {
			// <a class="simple-but border" href="polygon?2-1.ILinkListener-buffs-0-buff-buyLink"><span><span>Активировать за
			lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
			// Здесь бывают накладки, когда усиление прошло раньше
			if len(lstLink) < 2 {
				return
			}
			стрРез = lstLink[1]
			lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
			strLink := "https://wartank.ru/" + lstLink[0]
			lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
			for фнУсиление() {
				полигон.Обновить()
				lstPolygon = полигон.СписПолучить()
				lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
				if len(lstLink) < 2 {
					break
				}
				стрРез = lstLink[1]
				lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
				strLink := "https://wartank.ru/" + lstLink[0]
				lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
				log.Printf("усилениеДобавить(): link=%v", strLink)
			}
		}
		{ // Узнать на сколько форсирована атака
			if ind >= 7 && ind-7 < len(lstPolygon) {
				strForce := lstPolygon[ind-7]
				lstForce := strings.Split(strForce, `<span class="green2">+`)
				if len(lstForce) >= 2 {
					strForce = lstForce[1]
					lstForce = strings.Split(strForce, ` на `)
					strForce = lstForce[0]
					iForce, err := strconv.Atoi(strForce)
					if err == nil {
						танкСтат.ФорсажОбнов("attack", iForce)
						полигон.ПродуктСейчас().ИмяУст(стрАтака)
					}
				}
			}
		}
	case стрБроня: // Усиливаем броню
		фнУсиление := func() bool {
			// <span class="green2">усиление брони</span><br/>
			for ind, стрРез = range lstPolygon {
				if strings.Contains(стрРез, `>усиление брони<`) {
					ind += 8
					if ind >= len(lstPolygon) {
						return false
					}
					стрРез = lstPolygon[ind]
					if стрРез == "" {
						return false
					}
					if strings.Contains(стрРез, "Активировать за") {
						return false
					}
					return true
				}
			}
			return false
		}
		полигон.Обновить()
		lstPolygon = полигон.СписПолучить()
		if фнУсиление() {
			// <a class="simple-but border" href="polygon?2-1.ILinkListener-buffs-0-buff-buyLink"><span><span>Активировать за
			lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
			if len(lstLink) < 2 {
				return
			}
			стрРез = lstLink[1]
			lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
			strLink := "https://wartank.ru/" + lstLink[0]
			lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
			for фнУсиление() {
				полигон.Обновить()
				lstPolygon = полигон.СписПолучить()
				lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
				// Здесь бывают накладки, когда усиление прошло раньше
				if len(lstLink) < 2 {
					return
				}
				стрРез = lstLink[1]
				lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
				strLink := "https://wartank.ru/" + lstLink[0]
				lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
				log.Printf("усилениеДобавить(): link=%v", strLink)
			}
		}

		{ // Узнать на сколько форсирована броня
			if ind >= 7 && ind-7 < len(lstPolygon) {
				strForce := lstPolygon[ind-7]
				lstForce := strings.Split(strForce, `<span class="green2">+`)
				if len(lstForce) >= 2 {
					strForce = lstForce[1]
					lstForce = strings.Split(strForce, ` на `)
					strForce = lstForce[0]
					iForce, err := strconv.Atoi(strForce)
					if err == nil {
						танкСтат.ФорсажОбнов(стрБроня, iForce)
						полигон.ПродуктСейчас().ИмяУст(стрБроня)
					}
				}
			}
		}
	case стрТочность: // Усиливаем точность
		фнУсиление := func() bool {
			// span class="green2">улучшение точности</span><br/>
			for ind, стрРез = range lstPolygon {
				if strings.Contains(стрРез, `>улучшение точности<`) {
					ind += 8
					if ind >= len(lstPolygon) {
						return false
					}
					стрРез = lstPolygon[ind]
					if стрРез == "" {
						return false
					}
					if strings.Contains(стрРез, "Активировать за") {
						return false
					}
					return true
				}
			}

			return false
		}
		полигон.Обновить()
		lstPolygon = полигон.СписПолучить()
		if фнУсиление() {
			// <a class="simple-but border" href="polygon?2-1.ILinkListener-buffs-0-buff-buyLink"><span><span>Активировать за
			lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
			// Здесь бывают накладки, когда усиление прошло раньше
			if len(lstLink) < 2 {
				return
			}
			стрРез = lstLink[1]
			lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
			strLink := "https://wartank.ru/" + lstLink[0]
			lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
			for фнУсиление() {
				полигон.Обновить()
				lstPolygon = полигон.СписПолучить()
				lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
				if len(lstLink) < 2 {
					break
				}
				стрРез = lstLink[1]
				lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
				strLink := "https://wartank.ru/" + lstLink[0]
				lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
				log.Printf("усилениеДобавить(): link=%v", strLink)
			}
		}

		{ // Узнать на сколько форсирована точность
			if ind >= 7 && ind-7 < len(lstPolygon) {
				strForce := lstPolygon[ind-7]
				lstForce := strings.Split(strForce, `<span class="green2">+`)
				if len(lstForce) >= 2 {
					strForce = lstForce[1]
					lstForce = strings.Split(strForce, ` на `)
					strForce = lstForce[0]
					iForce, err := strconv.Atoi(strForce)
					if err == nil {
						танкСтат.ФорсажОбнов(стрТочность, iForce)
						полигон.ПродуктСейчас().ИмяУст(стрТочность)
					}
				}
			}
		}
	case стрПрочность: // Усиливаем мощность
		фнУсиление := func() bool {
			// <span class="green2">увеличение прочности</span><br/>
			for ind, стрРез = range lstPolygon {
				if strings.Contains(стрРез, `>увеличение прочности<`) {
					ind += 8
					if ind >= len(lstPolygon) {
						return false
					}
					стрРез = lstPolygon[ind]
					if стрРез == "" {
						return false
					}
					if strings.Contains(стрРез, "Активировать за") {
						return false
					}
					return true
				}
			}
			return false
		}

		полигон.Обновить()
		lstPolygon = полигон.СписПолучить()
		_ = фнУсиление()
		// <a class="simple-but border" href="polygon?2-1.ILinkListener-buffs-0-buff-buyLink"><span><span>Активировать за
		lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
		if len(lstLink) < 2 {
			return
		}
		стрРез = lstLink[1]
		lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
		strLink := "https://wartank.ru/" + lstLink[0]
		lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
		for фнУсиление() {
			полигон.Обновить()
			lstPolygon = полигон.СписПолучить()
			lstLink := strings.Split(стрРез, `<a class="simple-but border" href="`)
			if len(lstLink) < 2 {
				break
			}
			стрРез = lstLink[1]
			lstLink = strings.Split(стрРез, `"><span><span>Получить бесплатно</span></span></a>`)
			strLink := "https://wartank.ru/" + lstLink[0]
			lstPolygon = полигон.Сеть().ВебВоркер().Получ(strLink)
		}

		{ // Узнать на сколько форсирована прочность
			if ind >= 7 && ind-7 < len(lstPolygon) {
				strForce := lstPolygon[ind-7]
				lstForce := strings.Split(strForce, `<span class="green2">+`)
				if len(lstForce) >= 2 {
					strForce = lstForce[1]
					lstForce = strings.Split(strForce, ` на `)
					strForce = lstForce[0]
					iForce, err := strconv.Atoi(strForce)
					if err == nil {
						танкСтат.ФорсажОбнов(стрПрочность, iForce)
						полигон.ПродуктСейчас().ИмяУст(стрПрочность)
					}
				}
			}
		}
	default: // Неизвестно что
		Hassert(false, "усилениеДобавить(): усиление(%v) неизвестно", strParam)
	}
}
